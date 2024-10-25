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
	err = conn.AutoMigrate(&models.Repository{}, &models.WorkflowRun{}, &models.WorkflowStatistics{}, &models.JobStatistics{})
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
		Name:        repo.Name,
		Owner:       repo.Owner,
		FullName:    repo.FullName,
		Description: repo.Description,
		Private:     repo.Private,
		Fork:        repo.Fork,
		CreatedAt:   repo.CreatedAt,
		UpdatedAt:   repo.UpdatedAt,
		PushedAt:    repo.PushedAt,
		Size:        repo.Size,
		StarCount:   repo.StarCount,
		Language:    repo.Language,
		HasIssues:   repo.HasIssues,
		HasProjects: repo.HasProjects,
		HasWiki:     repo.HasWiki,
	}
	return db.Conn.Create(repository).Error
}

func (db *Database) DeleteRepository(id int) error {
	return db.Conn.Delete(&models.Repository{}, id).Error
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

func (db *Database) SaveWorkflow(workflow *github.Workflow) error {
	workflowModel := models.Workflow{
		WorkflowID: workflow.GetID(),
		NodeID:     workflow.GetNodeID(),
		Name:       workflow.GetName(),
		Path:       workflow.GetPath(),
		State:      workflow.GetState(),
		CreatedAt:  workflow.GetCreatedAt().Time,
		UpdatedAt:  workflow.GetUpdatedAt().Time,
	}
	return db.Conn.Create(workflowModel).Error
}

func (db *Database) DeleteWorkflow(id int) error {
	return db.Conn.Delete(&models.Workflow{}, id).Error
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

func (db *Database) GetWorkflowJob(workflowID int) (*models.Job, error) {
	var job models.Job
	err := db.Conn.First(&job, workflowID).Error
	return &job, err
}

func (db *Database) GetWorkflowJobs(workflowID int) ([]models.Job, error) {
	var jobs []models.Job
	err := db.Conn.Where("workflow_id = ?", workflowID).Find(&jobs).Error
	return jobs, err
}

func (db *Database) SaveWorkflowJob(job *github.WorkflowJob) error {
	jobModel := models.Job{
		JobID:           job.GetID(),
		RunID:           job.GetRunID(),
		RunURL:          job.GetRunURL(),
		NodeID:          job.GetNodeID(),
		HeadSHA:         job.GetHeadSHA(),
		URL:             job.GetURL(),
		HTMLURL:         job.GetHTMLURL(),
		Status:          job.GetStatus(),
		Conclusion:      job.GetConclusion(),
		CreatedAt:       job.GetCreatedAt().Time,
		CompletedAt:     job.GetCompletedAt().Time,
		Name:            job.GetName(),
		Steps:           []models.TaskStep{},
		CheckRunURL:     job.GetCheckRunURL(),
		Labels:          job.Labels,
		RunnerID:        job.GetRunnerID(),
		RunnerName:      job.GetRunnerName(),
		RunnerGroupID:   job.GetRunnerGroupID(),
		RunnerGroupName: job.GetRunnerGroupName(),
		RunAttempt:      int(job.GetRunAttempt()),
		WorkflowName:    job.GetWorkflowName(),
	}
	return db.Conn.Create(jobModel).Error
}

func (db *Database) DeleteWorkflowJob(id int) error {
	return db.Conn.Delete(&models.Job{}, id).Error
}

func (db *Database) GetJobStatistics() ([]models.JobStatistics, error) {
	var stats []models.JobStatistics
	err := db.Conn.Find(&stats).Error
	return stats, err
}

func (db *Database) SaveJobStatistics(stats *models.JobStatistics) error {
	jobStatistics := models.JobStatistics{
		ID: stats.ID,
	}
	return db.Conn.Save(jobStatistics).Error
}

func (db *Database) DeleteJobStatistics(id int) error {
	return db.Conn.Delete(&models.JobStatistics{}, id).Error
}

func (db *Database) GetWorkflowStatistics() ([]models.WorkflowStatistics, error) {
	var stats []models.WorkflowStatistics
	err := db.Conn.Find(&stats).Error
	return stats, err
}

func (db *Database) SaveWorkflowStatistics(stats *models.WorkflowStatistics) error {
	workflowStatistics := models.WorkflowStatistics{
		ID: stats.ID,
	}
	return db.Conn.Save(workflowStatistics).Error
}

func (db *Database) DeleteWorkflowStatistics(id int) error {
	return db.Conn.Delete(&models.WorkflowStatistics{}, id).Error
}
