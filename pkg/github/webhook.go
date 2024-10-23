package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github"
	"github.com/mooshe3/github-actions-aggregator/pkg/db"
)

type WebhookHandler struct {
	db            *db.Database
	webhookSecret []byte
}

func NewWebhookHandler(db *db.Database, secret string) *WebhookHandler {
	return &WebhookHandler{
		db:            db,
		webhookSecret: []byte(secret),
	}
}

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
	case *github.WorkflowRunEvent:
		wh.handleWorkflowRunEvent(e)
	default:
		// Unsupported event type
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}

func (wh *WebhookHandler) verifySignature(signature string, payload []byte) bool {
	mac := hmac.New(sha256.New, wh.webhookSecret)
	mac.Write(payload)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (wh *WebhookHandler) handleWorkflowRunEvent(event *github.WorkflowRunEvent) {
	action := event.GetAction()
	run := event.GetWorkflowRun()

	switch action {
	case "completed":
		// Save or update the workflow run in the database
		err := wh.db.SaveWorkflowRun(run)
		if err != nil {
			// Log error
		}

		// Enqueue a job to aggregate data after a new run is saved
		wh.workerPool.JobQueue <- worker.Job{
			Type: "aggregate_data",
		}

	case "requested":
		// Handle other actions if needed
	}
}
