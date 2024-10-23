package github

import (
	"context"

	gh "github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Client struct {
	ghClient *gh.Client
	ctx      context.Context
}

func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)

	return &Client{
		ghClient: client,
		ctx:      ctx,
	}
}

func (c *Client) ListWorkflows(owner, repo string) ([]*gh.Workflow, error) {
	workflows, _, err := c.ghClient.Actions.ListWorkflows(c.ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}
	return workflows.Workflows, nil
}

func (c *Client) ListWorkflowRuns(owner, repo string, workflowID int64) ([]*gh.WorkflowRun, error) {
	runs, _, err := c.ghClient.Actions.ListWorkflowRunsByID(c.ctx, owner, repo, workflowID, nil)
	if err != nil {
		return nil, err
	}
	return runs.WorkflowRuns, nil
}

func (c *Client) GetWorkflowRun(owner, repo string, runID int64) (*gh.WorkflowRun, error) {
	run, _, err := c.ghClient.Actions.GetWorkflowRunByID(c.ctx, owner, repo, runID)
	if err != nil {
		return nil, err
	}
	return run, nil
}
