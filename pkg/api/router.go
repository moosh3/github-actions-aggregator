package api

import (
	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/auth"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/github"
)

func StartServer(cfg *config.Config, db *db.Database, githubClient *github.Client) {
	r := gin.Default()

	// Public routes
	r.GET("/login", auth.GitHubLogin)
	r.GET("/callback", auth.GitHubCallback)

	// Webhook route (exclude middleware that could interfere)
	webhookHandler := github.NewWebhookHandler(db, githubClient, cfg.Github.WebhookSecret)
	r.POST("/webhook", webhookHandler.HandleWebhook)

	// Protected routes
	protected := r.Group("/", auth.AuthMiddleware())
	{
		protected.GET("/workflows/:id/stats", GetWorkflowStats)
	}

	r.Run(":" + cfg.ServerPort)
}
