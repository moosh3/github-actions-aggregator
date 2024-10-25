package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"gorm.io/gorm"
)

// GetRepository returns a single repository by ID.
func GetRepository(c *gin.Context) {
	repoIdParam := c.Param("id")
	repoId, err := strconv.ParseInt(repoIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var repo models.Repository
	err = db.Where("id = ?", repoId).First(&repo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve repository"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

func GetRepositories(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var repos []models.Repository
	err := db.Find(&repos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve repositories"})
		return
	}

	c.JSON(http.StatusOK, repos)
}

// GetRepositoryWorkflows returns all workflows for a given repository.
func GetRepositoryWorkflows(c *gin.Context) {
	repoIdParam := c.Param("id")
	repoId, err := strconv.ParseInt(repoIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var workflows []models.Workflow
	err = db.Where("repository_id = ?", repoId).Find(&workflows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows"})
		return
	}

	c.JSON(http.StatusOK, workflows)
}

func GetWorkflow(c *gin.Context) {
	workflowIdParam := c.Param("id")
	workflowId, err := strconv.ParseInt(workflowIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var workflow models.Workflow
	err = db.Where("id = ?", workflowId).First(&workflow).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func GetWorkflowJobs(c *gin.Context) {
	workflowIdParam := c.Param("id")
	workflowId, err := strconv.ParseInt(workflowIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var jobs []models.Job
	err = db.Where("workflow_id = ?", workflowId).Find(&jobs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// GetWorkflowRuns returns all runs for a given workflow.
func GetWorkflowRuns(c *gin.Context) {
	workflowIdParam := c.Param("id")
	workflowId, err := strconv.ParseInt(workflowIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var runs []models.WorkflowRun
	err = db.Where("workflow_id = ?", workflowId).Find(&runs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow runs"})
		return
	}

	c.JSON(http.StatusOK, runs)
}

// GetWorkflowRun returns a single workflow run by ID.
func GetWorkflowRun(c *gin.Context) {
	runIdParam := c.Param("id")
	runId, err := strconv.ParseInt(runIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid run ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var run models.WorkflowRun
	err = db.Where("id = ?", runId).First(&run).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow run"})
		return
	}

	c.JSON(http.StatusOK, run)
}

// GetJob returns a single job by ID.
func GetJob(c *gin.Context) {
	jobIdParam := c.Param("id")
	jobId, err := strconv.ParseInt(jobIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var job models.Job
	err = db.Where("id = ?", jobId).First(&job).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// GetJobSteps returns all steps for a given job.
func GetJobSteps(c *gin.Context) {
	jobIdParam := c.Param("id")
	jobId, err := strconv.ParseInt(jobIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var steps []models.TaskStep
	err = db.Where("job_id = ?", jobId).Find(&steps).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job steps"})
		return
	}

	c.JSON(http.StatusOK, steps)
}

// GetJobStats returns statistics for a given job.
func GetJobStats(c *gin.Context) {
	jobIdParam := c.Param("id")
	jobId, err := strconv.ParseInt(jobIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var stats models.JobStatistics
	err = db.Where("job_id = ?", jobId).First(&stats).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetWorkflowStats returns statistics for a given workflow.
func GetWorkflowStats(c *gin.Context) {
	workflowIDParam := c.Param("id")
	startTimeParam := c.Query("start_time")
	endTimeParam := c.Query("end_time")

	// Convert workflowID to integer
	workflowID, err := strconv.ParseInt(workflowIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Default start and end times
	defaultStartTime := time.Now().AddDate(0, 0, -30) // Default to 30 days ago
	defaultEndTime := time.Now()

	// Parse start_time
	startTime, err := parseTimeParameter(startTimeParam, defaultStartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse end_time
	endTime, err := parseTimeParameter(endTimeParam, defaultEndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Access the database
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	// Check if the workflow exists
	var workflow models.Workflow
	err = db.First(&workflow, "id = ?", workflowID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow"})
		}
		return
	}

	// Query workflow runs
	var runs []models.WorkflowRun
	err = db.Where("workflow_id = ?", workflowID).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Find(&runs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow runs"})
		return
	}

	// Ensure startTime is before endTime
	if !startTime.Before(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time must be before end_time"})
		return
	}

	// Initialize counters
	totalRuns := len(runs)
	successCount := 0
	failureCount := 0
	cancelledCount := 0
	timedOutCount := 0
	actionRequiredCount := 0

	for _, run := range runs {
		switch run.Conclusion {
		case "success":
			successCount++
		case "failure":
			failureCount++
		case "cancelled":
			cancelledCount++
		case "timed_out":
			timedOutCount++
		case "action_required":
			actionRequiredCount++
		}
	}

	// Calculate percentages
	var successRate, failureRate, cancelledRate, timedOutRate, actionRequiredRate float64
	if totalRuns > 0 {
		successRate = float64(successCount) / float64(totalRuns) * 100
		failureRate = float64(failureCount) / float64(totalRuns) * 100
		cancelledRate = float64(cancelledCount) / float64(totalRuns) * 100
		timedOutRate = float64(timedOutCount) / float64(totalRuns) * 100
		actionRequiredRate = float64(actionRequiredCount) / float64(totalRuns) * 100
	}

	// Respond with extended statistics
	c.JSON(http.StatusOK, gin.H{
		"workflow_id":           workflowID,
		"workflow_name":         workflow.Name,
		"total_runs":            totalRuns,
		"success_count":         successCount,
		"failure_count":         failureCount,
		"cancelled_count":       cancelledCount,
		"timed_out_count":       timedOutCount,
		"action_required_count": actionRequiredCount,
		"success_rate":          successRate,
		"failure_rate":          failureRate,
		"cancelled_rate":        cancelledRate,
		"timed_out_rate":        timedOutRate,
		"action_required_rate":  actionRequiredRate,
		"start_time":            startTime.Format(time.RFC3339),
		"end_time":              endTime.Format(time.RFC3339),
	})
}
