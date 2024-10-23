package worker

import (
	"time"
)

// Scheduler handles periodic scheduling of jobs.
type Scheduler struct {
	wp       *WorkerPool
	ticker   *time.Ticker
	stopChan chan struct{}
}

// NewScheduler initializes a new Scheduler.
func NewScheduler(wp *WorkerPool) *Scheduler {
	return &Scheduler{
		wp:       wp,
		stopChan: make(chan struct{}),
	}
}

// Start begins the scheduler.
func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(10 * time.Minute) // Set interval as needed
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.scheduleJobs()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// scheduleJobs enqueues jobs to the worker pool.
func (s *Scheduler) scheduleJobs() {
	// Enqueue an aggregate_data job
	s.wp.JobQueue <- Job{
		Type: "aggregate_data",
	}

	// Enqueue other periodic jobs as needed
}
