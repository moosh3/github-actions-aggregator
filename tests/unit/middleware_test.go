package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(auth.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
