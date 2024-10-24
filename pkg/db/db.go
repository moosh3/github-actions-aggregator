package db

import (
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Conn *gorm.DB
}

func InitDB(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=github_actions_aggregator sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	err = conn.AutoMigrate(&models.Repository{}, &models.WorkflowRun{}, &models.Statistics{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate schema: %w", err)
	}

	return &Database{Conn: conn}, nil
}

func (db *Database) GetRepository() (models.Repository, error) {
	var repo models.Repository
	err := db.Conn.Find(&repo).Error
	return repo, err
}

func (db *Database) GetRepositories() ([]models.Repository, error) {
	var repos []models.Repository
	err := db.Conn.Where("monitor = ?", true).Find(&repos).Error
	return repos, err
}

func (db *Database) SaveRepository(repo *models.Repository) error {
	repository := models.Repository{
		Name:  repo.Name,
		Owner: repo.Owner,
	}
	return db.Conn.Create(repository).Error
}

func (db *Database) DeleteRepository(id int) error {
	return db.Conn.Delete(&models.Repository{}, id).Error
}

func (db *Database) GetWorkflowRun(id int) (*models.WorkflowRun, error) {
	var run models.WorkflowRun
	err := db.Conn.First(&run, id).Error
	return &run, err
}

func (db *Database) GetWorkflowRuns(repoID int) ([]models.WorkflowRun, error) {
	var runs []models.WorkflowRun
	err := db.Conn.Where("repository_id = ?", repoID).Find(&runs).Error
	return runs, err
}

func (db *Database) SaveWorkflowRun(run *github.WorkflowRun) error {
	workflowRun := models.WorkflowRun{
		WorkflowID:   run.GetWorkflowID(),
		RepositoryID: run.GetRepository().GetID(),
		Status:       run.GetStatus(),
		Conclusion:   run.GetConclusion(),
		RunNumber:    run.GetRunNumber(),
		Event:        run.GetEvent(),
		CreatedAt:    run.GetCreatedAt().Time,
		UpdatedAt:    run.GetUpdatedAt().Time,
	}

	// Upsert operation
	return db.Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&workflowRun).Error
}

func (db *Database) DeleteWorkflowRun(id int) error {
	return db.Conn.Delete(&models.WorkflowRun{}, id).Error
}

func (db *Database) GetWorkflow(id int) (*models.Workflow, error) {
	var workflow models.Workflow
	err := db.Conn.First(&workflow, id).Error
	return &workflow, err
}

func (db *Database) GetWorkflows() ([]models.Workflow, error) {
	var workflows []models.Workflow
	err := db.Conn.Find(&workflows).Error
	return workflows, err
}

func (db *Database) SaveWorkflow(workflow *models.Workflow) error {
	workflowModel := models.Workflow{
		Name: workflow.Name,
	}
	return db.Conn.Create(workflowModel).Error
}

func (db *Database) DeleteWorkflow(id int) error {
	return db.Conn.Delete(&models.Workflow{}, id).Error
}

func (db *Database) GetStatistics() ([]models.Statistics, error) {
	var stats []models.Statistics
	err := db.Conn.Find(&stats).Error
	return stats, err
}

func (db *Database) SaveStatistics(stats *models.Statistics) error {
	workflowStatistics := models.Statistics{
		ID: stats.ID,
	}
	return db.Conn.Save(workflowStatistics).Error
}

func (db *Database) DeleteStatistics(id int) error {
	return db.Conn.Delete(&models.Statistics{}, id).Error
}
