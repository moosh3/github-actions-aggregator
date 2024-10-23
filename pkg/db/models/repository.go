package models

import (
	"time"

	"gorm.io/gorm"
)

// Repository represents a GitHub repository
type Repository struct {
	gorm.Model
	Name        string `gorm:"index;not null"`
	FullName    string `gorm:"uniqueIndex;not null"`
	Description string
	Private     bool
	Fork        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PushedAt    time.Time
	Size        int
	StarCount   int
	Language    string
	HasIssues   bool
	HasProjects bool
	HasWiki     bool
	OwnerID     uint
	Owner       GitHubUser `gorm:"foreignKey:OwnerID"`
}
