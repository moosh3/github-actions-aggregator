package integration_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/api"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAPIIntegration(t *testing.T) {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize the test database
	db, err := InitTestDB(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer db.Close()

	// Start the API server in test mode
	router := gin.Default()
	api.SetupRoutes(router, cfg, db)

	// Make a test request
	resp, err := http.Get("http://localhost:" + cfg.ServerPort + "/stats")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Optionally, parse and verify the response body
}

func InitTestDB(cfg *config.Config) (*gorm.DB, error) {
	// Logic to initialize a test database
}
