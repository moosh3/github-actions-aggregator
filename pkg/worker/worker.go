package worker

import (
	"log"
	"sync"

	"github.com/moosh3/github-actions-aggregator/pkg/db"
)

// Job represents a unit of work to be processed by a worker.
type Job struct {
	Type    string
	Payload interface{}
}

// WorkerPool manages a pool of workers to process jobs.
type WorkerPool struct {
	JobQueue   chan Job
	NumWorkers int
	db         *db.Database
	stopChan   chan struct{}
	wg         sync.WaitGroup
	scheduler  *Scheduler
}

// NewWorkerPool initializes a new WorkerPool.
func NewWorkerPool(db *db.Database, numWorkers int) *WorkerPool {
	return &WorkerPool{
		JobQueue:   make(chan Job, 100), // Adjust the buffer size as needed
		NumWorkers: numWorkers,
		db:         db,
		stopChan:   make(chan struct{}),
	}
}

// Start initializes the worker pool and starts processing jobs.
func (wp *WorkerPool) Start() {
	log.Printf("Starting worker pool with %d workers", wp.NumWorkers)
	for i := 0; i < wp.NumWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// Start the scheduler if needed
	wp.scheduler = NewScheduler(wp)
	wp.scheduler.Start()
}

// Stop gracefully shuts down the worker pool.
func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
	wp.wg.Wait()
	wp.scheduler.Stop()
}

// worker is a goroutine that processes jobs from the JobQueue.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started", id)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker %d recovered from panic: %v", id, r)
		}
	}()

	log.Printf("Worker %d started", id)

	for {
		select {
		case job := <-wp.JobQueue:

			// ...
		}
	}
}
