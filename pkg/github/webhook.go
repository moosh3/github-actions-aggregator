package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github"
	"github.com/moosh3/github-actions-aggregator/pkg/db"
	"github.com/moosh3/github-actions-aggregator/pkg/worker"
)

// WebhookHandler handles GitHub webhook events.
type WebhookHandler struct {
	db       *db.Database
	client   *Client
	whSecret []byte
	worker   *worker.WorkerPool
}

// NewWebhookHandler creates a new WebhookHandler instance.
//
// Parameters:
//   - db: A pointer to the database instance.
//   - client: A pointer to the GitHub client.
//   - secret: The webhook secret used for signature verification.
//   - worker: A pointer to the worker pool.
//
// Returns:
//   - A pointer to the new WebhookHandler instance.
func NewWebhookHandler(db *db.Database, client *Client, secret string, worker *worker.WorkerPool) *WebhookHandler {
	return &WebhookHandler{
		db:       db,
		client:   client,
		whSecret: []byte(secret),
		worker:   worker,
	}
}

// HandleWebhook processes incoming GitHub webhook events.
//
// It verifies the webhook signature, parses the event, and handles
// different event types accordingly.
//
// Parameters:
//   - c: The Gin context for the HTTP request.
func (wh *WebhookHandler) HandleWebhook(c *gin.Context) {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not read request body"})
		return
	}

	// Verify the signature
	signature := c.GetHeader("X-Hub-Signature-256")
	if !wh.verifySignature(signature, payload) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Parse the event
	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not parse webhook"})
		return
	}

	// Handle different event types
	switch e := event.(type) {
	case *github.WorkflowRunEvent: // WorkflowRunEvent is triggered when a GitHub Actions workflow run is requested or completed.
		wh.handleWorkflowRunEvent(e)
	case *github.WorkflowJobEvent: // WorkflowJobEvent is triggered when a job is queued, started or completed.
		wh.handleWorkflowJobEvent(e)

	default:
		// Unsupported event type
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}

// verifySignature checks if the provided signature matches the expected signature
// calculated from the payload and the webhook secret.
//
// Parameters:
//   - signature: The signature provided in the webhook header.
//   - payload: The raw payload of the webhook.
//
// Returns:
//   - A boolean indicating whether the signature is valid.
func (wh *WebhookHandler) verifySignature(signature string, payload []byte) bool {
	mac := hmac.New(sha256.New, wh.whSecret)
	mac.Write(payload)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// handleWorkflowRunEvent processes GitHub workflow run events.
//
// It handles different actions for workflow runs, such as saving completed
// runs to the database and enqueueing data aggregation jobs.
//
// Parameters:
//   - event: A pointer to the GitHub WorkflowRunEvent.
func (wh *WebhookHandler) handleWorkflowRunEvent(event *github.WorkflowRunEvent) {
	action := event.GetAction()
	workflow := event.GetWorkflow()
	run := event.GetWorkflowRun()

	switch action {
	case "completed":
		err := wh.db.SaveWorkflow(workflow)
		if err != nil {
			// Log error
		}
		// Save or update the workflow run in the database
		err = wh.db.SaveWorkflowRun(run)
		if err != nil {
			// Log error
		}

		// Enqueue a job to aggregate data after a new run is saved
		wh.worker.JobQueue <- worker.Job{
			Type: "aggregate_data",
		}

	case "requested":
		// Handle other actions if needed
	}
}

func (wh *WebhookHandler) handleWorkflowJobEvent(event *github.WorkflowJobEvent) {
	job := event.GetWorkflowJob()
	err := wh.db.SaveWorkflowJob(job)
	if err != nil {
		// Log error
	}
}
