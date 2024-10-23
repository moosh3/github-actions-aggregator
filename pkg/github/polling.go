package github

import (
	"log"
	"time"

	gh "github.com/google/go-github/v50/github"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
)

type Poller struct {
	db       *db.Database
	client   *Client
	interval time.Duration
}

func NewPoller(db *db.Database, client *Client, interval time.Duration) *Poller {
	return &Poller{
		db:       db,
		client:   client,
		interval: interval,
	}
}

func (p *Poller) Start() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.pollRepositories()
		}
	}
}

func (p *Poller) pollRepositories() {
	repos, err := p.db.GetMonitoredRepositories()
	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
		return
	}

	for _, repo := range repos {
		p.pollWorkflows(repo)
	}
}

func (p *Poller) pollWorkflows(repo models.Repository) {
	workflows, err := p.client.ListWorkflows(repo.Owner, repo.Name)
	if err != nil {
		log.Printf("Error listing workflows for %s/%s: %v", repo.Owner, repo.Name, err)
		return
	}

	for _, workflow := range workflows {
		p.pollWorkflowRuns(repo.Owner, repo.Name, workflow)
	}
}

func (p *Poller) pollWorkflowRuns(owner, repoName string, workflow *gh.Workflow) {
	// Use client.go methods
	runs, err := p.client.ListWorkflowRuns(owner, repoName, *workflow.ID)
	if err != nil {
		log.Printf("Error listing workflow runs: %v", err)
		return
	}

	// Process runs...
}
