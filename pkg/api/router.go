package api

import (
	"github.com/gin-gonic/gin"
	"github.com/moosh3/github-actions-aggregator/pkg/auth"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/github"
	"github.com/moosh3/github-actions-aggregator/pkg/worker"
)

func StartServer(cfg *config.Config, db *db.Database, githubClient *github.Client, worker *worker.WorkerPool) {
	r := gin.Default()

	// Enable CORS
	r.Use(corsMiddleware())

	// Public routes for Github OAuth
	r.GET("/login", auth.GitHubLogin)
	r.GET("/callback", auth.GitHubCallback)

	// Webhook route for Github events (exclude middleware that could interfere)
	webhookHandler := github.NewWebhookHandler(db, githubClient, cfg.GitHub.WebhookSecret, worker)
	r.POST("/webhook", webhookHandler.HandleWebhook)

	auth := r.Group("/auth")
	{
		auth.GET("/github/login", handleGithubLogin(cfg))
		auth.GET("/github/callback", handleGithubCallback(cfg))
		auth.GET("/user", authMiddleware(cfg), getCurrentUser)
	}

	// Require authentication for all repository routes
	protected := r.Group("/repositories", authMiddleware(cfg))
	{
		protected.GET("", GetRepositories)
		protected.GET("/:repoId", GetRepository)
		protected.GET("/:repoId/workflows", GetRepositoryWorkflows)                    // Get all workflows for a repository
		protected.GET("/:repoId/workflows/:workflowId", GetWorkflow)                   // Get a specific workflow
		protected.GET("/:repoId/workflows/:workflowId/runs", GetWorkflowRuns)          // Get all runs for a workflow
		protected.GET("/:repoId/workflows/:workflowId/runs/:runId", GetWorkflowRun)    // Get a specific run
		protected.GET("/:repoId/workflows/:workflowId/stats", GetWorkflowStats)        // Get stats for a workflow
		protected.GET("/:repoId/workflows/:workflowId/jobs", GetWorkflowJobs)          // Get all jobs for a workflow
		protected.GET("/:repoId/workflows/:workflowId/jobs/:jobId", GetJob)            // Get a specific job
		protected.GET("/:repoId/workflows/:workflowId/jobs/:jobId/steps", GetJobSteps) // Get all steps for a job
		protected.GET("/:repoId/workflows/:workflowId/jobs/:jobId/stats", GetJobStats) // Get stats for a job
	}

	r.Run(":" + cfg.ServerPort)
}
