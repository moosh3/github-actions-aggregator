package models

import "time"

type Statistics struct {
	ID          uint `gorm:"primaryKey"`
	TotalRuns   int64
	SuccessRate float64
	// Add fields for additional statistics
	UpdatedAt time.Time
}
