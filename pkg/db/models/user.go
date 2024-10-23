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
	GravatarID        string
	URL               string
	HTMLURL           string
	FollowersURL      string
	FollowingURL      string
	GistsURL          string
	StarredURL        string
	SubscriptionsURL  string
	OrganizationsURL  string
	ReposURL          string
	EventsURL         string
	ReceivedEventsURL string
	Type              string
	SiteAdmin         bool
	Name              string
	Company           string
	Blog              string
	Location          string
	Email             string
	Hireable          bool
	Bio               string
	TwitterUsername   string
	PublicRepos       int
	PublicGists       int
	Followers         int
	Following         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
