package models_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"github.com/stretchr/testify/assert"
)

func NewTestUser() *models.User {
	return &models.User{
		Username: fmt.Sprintf("user_%d", rand.Int()),
		Email:    fmt.Sprintf("user_%d@example.com", rand.Int()),
	}
}

func TestUserModel(t *testing.T) {
	user := NewTestUser()

	err := user.Validate()
	assert.NoError(t, err)

	// Test other methods of the User model
}
