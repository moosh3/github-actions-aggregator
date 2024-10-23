package auth_test

import (
	"testing"

	"github.com/moosh3/github-actions-aggregator/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestGenerateStateToken(t *testing.T) {
	token1 := auth.GenerateStateToken()
	token2 := auth.GenerateStateToken()
	assert.NotEqual(t, token1, token2)
	assert.Len(t, token1, 32)
}
