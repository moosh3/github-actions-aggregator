package integration_test

import (
	"testing"

	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseIntegration(t *testing.T) {
	// Initialize test database
	cfg := config.LoadTestConfig()
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer dbConn.Close()

	// Run migrations or use a pre-migrated test database
	dbConn.AutoMigrate(&models.User{})

	// Perform database operations
	user := &models.User{
		Username: "integrationuser",
		Email:    "integration@example.com",
	}

	err = dbConn.Create(user).Error
	assert.NoError(t, err)

	// Retrieve and verify the user
	var retrievedUser models.User
	err = dbConn.First(&retrievedUser, "username = ?", "integrationuser").Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, retrievedUser.Email)
}
