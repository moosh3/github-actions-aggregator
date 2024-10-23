package github

type GitHubClient interface {
	GetWorkflowRuns(repo string) ([]WorkflowRun, error)
}
