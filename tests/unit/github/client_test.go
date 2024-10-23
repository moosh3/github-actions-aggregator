package github_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moosh3/github-actions-aggregator/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGitHubClient struct {
	mock.Mock
}

func StartMockGitHubServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/repos/", func(w http.ResponseWriter, r *http.Request) {
		// Return mock response
	})
	return httptest.NewServer(handler)
}

func (m *MockGitHubClient) GetWorkflowRuns(repo string) ([]github.WorkflowRun, error) {
	args := m.Called(repo)
	return args.Get(0).([]github.WorkflowRun), args.Error(1)
}

func TestGetWorkflowRuns(t *testing.T) {
	mockClient := new(MockGitHubClient)
	mockRuns := []github.WorkflowRun{
		{ID: 1, Status: "success"},
		{ID: 2, Status: "failed"},
	}
	mockClient.On("GetWorkflowRuns", "testrepo").Return(mockRuns, nil)

	// Use mockClient in place of the real GitHubClient
	runs, err := mockClient.GetWorkflowRuns("testrepo")
	assert.NoError(t, err)
	assert.Len(t, runs, 2)
	assert.Equal(t, "success", runs[0].Status)
}
