package github

import (
	"context"
	"log"
	"sync"
	"time"

	gh "github.com/google/go-github/v50/github"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/db/models"
	"golang.org/x/oauth2"
)

const maxConcurrentPolls = 10

type Poller struct {
	db       *db.Database
	ghClient *gh.Client
	interval time.Duration
}

func NewPoller(db *db.Database, token *oauth2.Token, interval time.Duration) *Poller {
	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(context.Background(), ts)
	client := gh.NewClient(tc)

	return &Poller{
		db:       db,
		ghClient: client,
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

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentPolls)

	for _, repo := range repos {
		wg.Add(1)
		sem <- struct{}{}
		go func(repo models.Repository) {
			defer wg.Done()
			p.pollWorkflows(repo)
			<-sem
		}(repo)
	}

	wg.Wait()
}

func (p *Poller) pollWorkflows(repo models.Repository) {
	owner := repo.Owner
	repoName := repo.Name

	// List workflows
	workflows, _, err := p.ghClient.Actions.ListWorkflows(context.Background(), owner, repoName, nil)
	if err != nil {
		log.Printf("Error listing workflows for %s/%s: %v", owner, repoName, err)
		return
	}

	for _, workflow := range workflows.Workflows {
		p.pollWorkflowRuns(owner, repoName, workflow)
	}
}

func (p *Poller) pollWorkflowRuns(owner string, repoName string, workflow *gh.Workflow) {
	opts := &gh.ListWorkflowRunsOptions{
		ListOptions: gh.ListOptions{PerPage: 50},
	}

	runs, _, err := p.ghClient.Actions.ListWorkflowRunsByID(context.Background(), owner, repoName, *workflow.ID, opts)
	if err != nil {
		log.Printf("Error listing workflow runs for %s/%s (Workflow ID: %d): %v", owner, repoName, *workflow.ID, err)
		return
	}

	for _, run := range runs.WorkflowRuns {
		// Save or update workflow run in the database
		err := p.db.SaveWorkflowRun(run)
		if err != nil {
			log.Printf("Error saving workflow run ID %d: %v", *run.ID, err)
		}
	}
}

func (p *Poller) handleRateLimit(resp *gh.Response) {
	if resp.Rate.Remaining == 0 {
		resetTime := time.Until(resp.Rate.Reset.Time)
		log.Printf("Rate limit exceeded. Waiting for %v", resetTime)
		time.Sleep(resetTime)
	}
}
