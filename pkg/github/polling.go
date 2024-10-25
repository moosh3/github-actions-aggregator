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

// Poller represents a GitHub poller that periodically fetches workflow information.
type Poller struct {
	db       *db.Database
	ghClient *gh.Client
	interval time.Duration
}

// NewPoller creates a new Poller instance with the given database, OAuth token, and polling interval.
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

// Start begins the polling process, periodically calling pollRepositories based on the set interval.
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

// pollRepositories fetches all repositories accessible to the authenticated user and polls their workflows concurrently.
func (p *Poller) pollRepositories() {
	opt := &gh.RepositoryListOptions{
		ListOptions: gh.ListOptions{PerPage: 10}, // Adjust per page as needed
	}

	var allRepos []*gh.Repository
	for {
		repos, resp, err := p.ghClient.Repositories.List(context.Background(), "", opt)
		if err != nil {
			log.Printf("Error fetching repositories: %v", err)
			return
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentPolls)

	for _, repo := range allRepos {
		wg.Add(1)
		sem <- struct{}{}
		go func(repo *gh.Repository) {
			defer wg.Done()
			dbRepo := models.Repository{
				Owner: models.GitHubUser{
					Email: func(email *string) string {
						if email != nil {
							return *email
						}
						return ""
					}(repo.Owner.Email),
				},
				Name: *repo.Name, // Dereference the pointer
			}
			err := p.db.SaveRepository(&dbRepo)
			if err != nil {
				log.Printf("Error saving repository %s: %v", repo.Name, err)
			}
			<-sem
		}(repo)
	}

	wg.Wait()
}

// pollWorkflows fetches and processes all workflows for a given repository.
func (p *Poller) pollWorkflows(repo models.Repository) {
	owner := string(repo.Owner.Email)
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

// pollWorkflowRuns fetches and saves the runs for a specific workflow.
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

// handleRateLimit checks the rate limit from the GitHub API response and waits if the limit is exceeded.
func (p *Poller) handleRateLimit(resp *gh.Response) {
	if resp.Rate.Remaining == 0 {
		resetTime := time.Until(resp.Rate.Reset.Time)
		log.Printf("Rate limit exceeded. Waiting for %v", resetTime)
		time.Sleep(resetTime)
	}
}
