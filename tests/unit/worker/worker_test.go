package worker_test

import (
	"testing"

	"github.com/moosh3/github-actions-aggregator/pkg/worker"
	"github.com/stretchr/testify/assert"
)

func TestProcessJob(t *testing.T) {
	// Mock dependencies (e.g., database)
	db := &MockDatabase{}
	wp := worker.NewWorkerPool(db, 1)

	job := worker.Job{
		Type: "aggregate_data",
	}

	wp.processJob(job)

	// Assert that the aggregation function was called
	assert.True(t, db.AggregateDataCalled)
}
