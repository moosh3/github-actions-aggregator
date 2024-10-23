package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	// Placeholder for actual logic
	stats := map[string]interface{}{
		"total_runs":   100,
		"success_rate": 95.0,
	}
	c.JSON(http.StatusOK, stats)
}
