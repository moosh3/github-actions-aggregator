package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/moosh3/github-actions-aggregator/pkg/api"
	"github.com/moosh3/github-actions-aggregator/pkg/config"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/github"
	"github.com/moosh3/github-actions-aggregator/pkg/logger"
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

	// Start the API server
	api.StartServer(cfg, database, githubClient)
}

func runMigrations() error {
	cmd := exec.Command("./scripts/migrate.sh", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
