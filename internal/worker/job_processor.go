package worker

import (
	"context"
	"sync"
	"time"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/rs/zerolog"
)

type Job struct {
	Incident *model.Incident
	Logger   zerolog.Logger
}

type JobProcessor struct {
	JobQueue    chan Job
	Logger      zerolog.Logger
	workerCount int
	wg          sync.WaitGroup
}

const defaultQueueSize = 1000 // size for backpressure
const defaultWorkers = 5      // size of work pool

func NewJobProcessor(logger zerolog.Logger) *JobProcessor {
	return &JobProcessor{
		JobQueue:    make(chan Job, defaultQueueSize), // initializes the buffered channel
		Logger:      logger,
		workerCount: defaultWorkers,
	}
}

// Start launches the worker goroutines
func (p *JobProcessor) Start(ctx context.Context) {
	p.Logger.Info().Int("workers", p.workerCount).Int("queue_size", cap(p.JobQueue)).Msg("Starting in-process worker pool")

	for i := 1; i <= p.workerCount; i++ {
		p.wg.Add(1) // add a count for each worker we start
		go p.startWorker(ctx, i)
	}
}

func (p *JobProcessor) startWorker(ctx context.Context, id int) {
	defer p.wg.Done() // decrement the counter when the worker exits
	workerLog := p.Logger.With().Int("worker_id", id).Logger()

	// this loop runs forever until the service shuts down
	for {
		select {
		case <-ctx.Done():
			// GRACEFUL SHUTDOWN: Context was cancelled (e.g., Ctrl+C), worker finishes its work and exits.
			workerLog.Info().Msg("Worker shutting down gracefully.")
			return
		case job, ok := <-p.JobQueue:
			if !ok {
				// channel was closed (via Stop()), worker exits.
				workerLog.Info().Msg("Job queue closed, worker exiting.")
				return
			}
			p.process(job)
		}
	}
}

// process simulates the slow, real-world I/O notification task.
func (p *JobProcessor) process(job Job) {
	// LOGGING ASYNC: We use the logger passed in the Job to maintain traceability
	log := job.Logger.With().Str("task_id", job.Incident.ID).Logger()
	log.Info().Msg("Worker received job. Simulating slow external notification I/O...")

	// Simulate the slow network call that requires a Goroutine
	time.Sleep(3 * time.Second)

	log.Info().Msgf("Notification for Incident %s completed successfully.", job.Incident.ID)
}

// Stop closes the job channel, signaling workers to stop processing new jobs.
func (p *JobProcessor) Stop() {
	close(p.JobQueue)
	p.Logger.Info().Msg("Job queue closed.")
}

// Wait blocks until all workers have called Done().
func (p *JobProcessor) Wait() {
	p.wg.Wait()
}
