package models

import (
	"time"

	"gorm.io/gorm"
)

// GitHubUser represents a GitHub user based on the GitHub API
type GitHubUser struct {
	gorm.Model
	Login             string `gorm:"uniqueIndex;not null"`
	ID                int64  `gorm:"uniqueIndex;not null"`
	NodeID            string `gorm:"not null"`
	AvatarURL         string
	TenantID          string // The tenant ID for the user
	URL               string
	HTMLURL           string
	SubscriptionsURL  string
	OrganizationsURL  string
	ReposURL          string
	EventsURL         string
	ReceivedEventsURL string
	Type              string
	SiteAdmin         bool
	Name              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
