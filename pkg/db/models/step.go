package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStep struct {
	gorm.Model
	Name        string
	Status      string
	Conclusion  string
	StartedAt   time.Time
	CompletedAt time.Time
}
