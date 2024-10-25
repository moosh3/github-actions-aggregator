package models

import "time"

type WorkflowStatistics struct {
	ID          uint `gorm:"primaryKey"`
	TotalRuns   int64
	SuccessRate float64
	// Add fields for additional statistics
	UpdatedAt time.Time
}

type JobStatistics struct {
	ID        uint `gorm:"primaryKey"`
	UpdatedAt time.Time
}
