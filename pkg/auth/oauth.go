package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		Endpoint:     github.Endpoint,
		RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/callback"),
		Scopes:       []string{"read:user", "repo", "workflow"},
	}
)

func GitHubLogin(c *gin.Context) {
	state := generateStateString()
	saveOAuthState(c, state)
	url := oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func GitHubCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	savedState := getSavedOAuthState(c)
	if state != savedState {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	user, err := getUserInfo(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Save or update user in database (implement saveOrUpdateUser)
	// Set user session (implement setUserSession)

	c.Redirect(http.StatusFound, "/dashboard")
}

func getUserInfo(token *oauth2.Token) (*GitHubUser, error) {
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func generateStateString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func saveOAuthState(c *gin.Context, state string) {
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
}

func getSavedOAuthState(c *gin.Context) string {
	state, err := c.Cookie("oauth_state")
	if err != nil {
		return ""
	}
	return state
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
