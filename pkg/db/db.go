package db

import (
	"github.com/google/go-github/v50/github"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Conn *gorm.DB
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Database initialization logic
}

func (db *Database) GetMonitoredRepositories() ([]models.Repository, error) {
	var repos []models.Repository
	err := db.Conn.Where("monitor = ?", true).Find(&repos).Error
	return repos, err
}

func (db *Database) SaveWorkflowRun(run *github.WorkflowRun) error {
	workflowRun := models.WorkflowRun{
		ID:           run.GetID(),
		WorkflowID:   run.GetWorkflowID(),
		RepositoryID: run.GetRepository().GetID(),
		Status:       run.GetStatus(),
		Conclusion:   run.GetConclusion(),
		RunNumber:    run.GetRunNumber(),
		Event:        run.GetEvent(),
		CreatedAt:    run.GetCreatedAt().Time,
		UpdatedAt:    run.GetUpdatedAt().Time,
		// Add other fields as needed
	}

	// Upsert operation
	return db.Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&workflowRun).Error
}

func (db *Database) SaveStatistics(stats *models.Statistics) error {
	// Upsert operation
	return db.Conn.Save(stats).Error
}
