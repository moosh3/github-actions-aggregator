package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	ID              int64
	JobID           int64
	RunID           int64
	RunURL          string
	NodeID          string
	HeadSHA         string
	URL             string
	HTMLURL         string
	LogsURL         string
	CheckRunURL     string
	RunnerID        int64
	CreatedAt       time.Time
	Name            string
	Labels          []string
	RunAttempt      int
	RunnerName      string
	RunnerGroupID   int64
	RunnerGroupName string
	WorkflowID      int64
	WorkflowName    string
	Status          string
	Conclusion      string
	CompletedAt     time.Time
	Steps           []TaskStep
}
