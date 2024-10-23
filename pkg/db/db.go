package db

import (
	"fmt"

	"github.com/google/go-github/v50/github"
	config "github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Conn *gorm.DB
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Database initialization logic
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Repository{}, &models.WorkflowRun{}, &models.Statistics{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}

	return db, nil
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
