package models

import (
	"time"

	"gorm.io/gorm"
)

// WorkflowRun represents a workflow run from GitHub API
type WorkflowRun struct {
	gorm.Model
	WorkflowID       int64 `gorm:"index"`
	Name             string
	HeadBranch       string
	HeadSHA          string
	Status           string
	Conclusion       string
	EventType        string
	URL              string
	HTMLURL          string
	JobsURL          string
	LogsURL          string
	CheckSuiteURL    string
	ArtifactsURL     string
	CancelURL        string
	RerunURL         string
	WorkflowURL      string
	RunNumber        int
	RunAttempt       int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	RunStartedAt     *time.Time
	JobsCount        int
	PullRequests     []PullRequest `gorm:"many2many:workflow_run_pull_requests;"`
	RepositoryID     int64         `gorm:"index"`
	HeadRepository   Repository    `gorm:"foreignKey:HeadRepositoryID"`
	HeadRepositoryID int64
}

// PullRequest represents a pull request associated with a workflow run
type PullRequest struct {
	gorm.Model
	URL    string
	Number int
}
