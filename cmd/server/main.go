package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/moosh3/github-actions-aggregator/pkg/github"
	"github.com/mooshe3/github-actions-aggregator/pkg/api"
	"github.com/mooshe3/github-actions-aggregator/pkg/config"
	"github.com/mooshe3/github-actions-aggregator/pkg/db"
	"github.com/mooshe3/github-actions-aggregator/pkg/logger"
	"github.com/mooshe3/github-actions-aggregator/pkg/worker"
)

func main() {
	// Initialize configurations
	cfg := config.LoadConfig()

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize GitHub client
	githubClient := github.NewClient(cfg.GitHub.AccessToken)

	// Start the poller if enabled
	if cfg.EnablePolling {
		poller := github.NewPoller(database, githubClient, cfg.PollingInterval)
		go poller.Start()
	}

	// Start the worker pool
	wp := worker.NewWorkerPool(database, 5) // Adjust the number of workers as needed
	wp.Start()

	// Start the API server
	go api.StartServer(cfg, database, githubClient)

	// Wait for interrupt signal to gracefully shut down the worker pool
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker pool...")
	wp.Stop()

	log.Println("Server exiting")
}
