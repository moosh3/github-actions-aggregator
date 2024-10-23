// pkg/api/handlers.go

package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"gorm.io/gorm"
)

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

	// Parse start and end times
	var startTime, endTime time.Time
	if startTimeParam != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endTimeParam != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	if !startTime.Before(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time must be before end_time"})
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

	// Calculate statistics
	totalRuns := len(runs)
	successCount := 0
	failureCount := 0

	for _, run := range runs {
		switch run.Conclusion {
		case "success":
			successCount++
		case "failure":
			failureCount++
		}
	}

	// Calculate percentages
	var successRate, failureRate float64
	if totalRuns > 0 {
		successRate = float64(successCount) / float64(totalRuns) * 100
		failureRate = float64(failureCount) / float64(totalRuns) * 100
	}

	// Respond with statistics
	c.JSON(http.StatusOK, gin.H{
		"workflow_id":   workflowID,
		"workflow_name": workflow.Name,
		"total_runs":    totalRuns,
		"success_count": successCount,
		"failure_count": failureCount,
		"success_rate":  successRate,
		"failure_rate":  failureRate,
		"start_time":    startTime.Format(time.RFC3339),
		"end_time":      endTime.Format(time.RFC3339),
	})
}
