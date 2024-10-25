package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func githubOAuthConfig(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GitHub.ClientID,
		ClientSecret: config.GitHub.ClientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func handleGithubLogin(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthConfig := githubOAuthConfig(config)
		url := oauthConfig.AuthCodeURL("state")
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func handleGithubCallback(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")

		// Exchange code for token
		oauthConfig := githubOAuthConfig(config)
		token, err := oauthConfig.Exchange(c, code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
			return
		}

		// Get user info from GitHub
		client := oauthConfig.Client(c, token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
			return
		}
		defer resp.Body.Close()

		var githubUser struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
			Login string `json:"login"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
			return
		}

		// Create JWT
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId":   githubUser.ID,
			"email":    githubUser.Email,
			"username": githubUser.Login,
		})

		tokenString, err := jwtToken.SignedString([]byte(config.GitHub.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		// Redirect to frontend with token
		c.Redirect(http.StatusTemporaryRedirect,
			config.FrontendURL+"/auth/callback?token="+tokenString)
	}
}
