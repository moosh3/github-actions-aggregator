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
