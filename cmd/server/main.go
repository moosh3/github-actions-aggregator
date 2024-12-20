package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/moosh3/github-actions-aggregator/pkg/api"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/github"
	"github.com/moosh3/github-actions-aggregator/pkg/logger"
	"github.com/moosh3/github-actions-aggregator/pkg/worker"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Run migrations
	err := runMigrations()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize GitHub client
	githubClient := github.NewClient(cfg.GitHub.AccessToken)

	// Initialize worker pool for polling
	pollingWorkerPool := worker.NewWorkerPool(database, cfg.PollingWorkerPoolSize)
	pollingWorkerPool.Start()

	// Initialize worker pool for webhooks
	webhookWorkerPool := worker.NewWorkerPool(database, cfg.WebhookWorkerPoolSize)
	webhookWorkerPool.Start()

	// Start the API server
	go api.StartServer(cfg, database, githubClient, webhookWorkerPool)

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop the worker pools
	webhookWorkerPool.Stop()
	pollingWorkerPool.Stop()
	log.Println("Server exiting")
}

func runMigrations() error {
	cmd := exec.Command("./scripts/migrate.sh", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
