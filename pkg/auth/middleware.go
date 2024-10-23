package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func getUserFromSession(c *gin.Context) *models.GitHubUser {
	// Retrieve user information from session or cookie
	// For example, using a secure cookie:
	userID, err := c.Cookie("user_id")
	if err != nil {
		return nil
	}
	// Convert userID to int64
	id, _ := strconv.ParseInt(userID, 10, 64)
	// Fetch user from database using userID
	// Return the user object
	return &models.GitHubUser{ID: id}
}
