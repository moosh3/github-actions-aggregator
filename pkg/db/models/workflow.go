package models

import (
	"time"

	"gorm.io/gorm"
)

// Workflow represents a GitHub workflow
type Workflow struct {
	gorm.Model
	Name         string    `gorm:"type:varchar(255);not null"`
	Path         string    `gorm:"type:varchar(255);not null"`
	State        string    `gorm:"type:varchar(50)"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	URL          string    `gorm:"type:varchar(255)"`
	HTMLURL      string    `gorm:"column:html_url;type:varchar(255)"`
	BadgeURL     string    `gorm:"column:badge_url;type:varchar(255)"`
	RepositoryID uint      `gorm:"not null"`
	// You might want to add a foreign key relationship to a Repository model if you have one
	// Repository   Repository `gorm:"foreignKey:RepositoryID"`
}

// TableName specifies the table name for the Workflow model
func (Workflow) TableName() string {
	return "workflows"
}
