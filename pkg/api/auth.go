package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuthState struct {
	State    string
	TenantID string
	ReturnTo string
}

// Redis client for storing state
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func getGithubOAuthConfig(cfg *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.GitHub.ClientID,
		ClientSecret: cfg.GitHub.ClientSecret,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/auth/github/callback", cfg.APIURL),
	}
}

func handleGithubLogin(c *gin.Context, cfg *config.Config) {
	// Generate random state
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state with tenant info in Redis
	oauthState := OAuthState{
		State:    state,
		ReturnTo: c.Query("returnTo"), // Optional return URL
	}

	stateJSON, err := json.Marshal(oauthState)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal state"})
		return
	}

	// Store state in Redis with 15-minute expiration
	err = redisClient.Set(context.Background(),
		fmt.Sprintf("oauth_state:%s", state),
		string(stateJSON),
		15*time.Minute,
	).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store state"})
		return
	}

	oauthConfig := getGithubOAuthConfig(cfg)

	// Generate authorization URL
	authURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
	)

	// Redirect to GitHub
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func handleGithubCallback(c *gin.Context, cfg *config.Config) {
	code := c.Query("code")
	state := c.Query("state")

	// Validate state
	stateKey := fmt.Sprintf("oauth_state:%s", state)
	stateJSON, err := redisClient.Get(context.Background(), stateKey).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired state"})
		return
	}

	// Parse stored state
	var oauthState OAuthState
	if err := json.Unmarshal([]byte(stateJSON), &oauthState); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state data"})
		return
	}

	// Delete used state
	redisClient.Del(context.Background(), stateKey)

	oauthConfig := getGithubOAuthConfig(cfg)

	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Get GitHub user info
	client := oauthConfig.Client(c, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var githubUser models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Create or update user in your database
	user, err := db.UpdateUser(models.GitHubUser{
		Email:     githubUser.Email,
		Username:  githubUser.Login,
		Name:      githubUser.Name,
		AvatarURL: githubUser.AvatarURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user"})
		return
	}

	// Generate JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := jwtToken.SignedString([]byte(cfg.GitHub.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Determine redirect URL
	redirectURL := fmt.Sprintf("https://%s/auth/callback?token=%s", tenant.Domain, tokenString)
	if oauthState.ReturnTo != "" {
		// Validate and sanitize ReturnTo URL here
		redirectURL = fmt.Sprintf("%s&returnTo=%s", redirectURL, oauthState.ReturnTo)
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
