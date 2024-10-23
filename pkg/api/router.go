package api

import (
	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/auth"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/github"
)

func StartServer(cfg *config.Config) {
	r := gin.Default()

	// Public routes
	r.GET("/login", auth.GitHubLogin)
	r.GET("/callback", auth.GitHubCallback)

	// Webhook route (exclude middleware that could interfere)
	webhookHandler := github.NewWebhookHandler(db, cfg.GitHub.WebhookSecret)
	r.POST("/webhook", webhookHandler.HandleWebhook)

	// Protected routes
	protected := r.Group("/", auth.AuthMiddleware())
	{
		protected.GET("/dashboard", dashboardHandler)
		protected.GET("/stats", statsHandler)
		// Add other protected routes
	}

	r.Run(":" + cfg.ServerPort)
}
