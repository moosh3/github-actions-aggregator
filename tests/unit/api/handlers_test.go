package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestGetStats(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new router without middleware
	router := gin.New()
	router.GET("/stats", api.GetStats)

	// Create a request to pass to our handler
	req, err := http.NewRequest(http.MethodGet, "/stats", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Optionally, check the response body
	expected := `{"total_runs":100,"success_rate":95}`
	assert.JSONEq(t, expected, w.Body.String())
}
