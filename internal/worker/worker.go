package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hascho/go-incident-dashboard-api/internal/repository"
	"github.com/rs/zerolog"
)

type NotificationWorker struct {
	Repo   repository.JobRepository
	Logger zerolog.Logger
	ID     string
}

func NewNotificationWorker(repo repository.JobRepository, logger zerolog.Logger, id string) *NotificationWorker {
	return &NotificationWorker{
		Repo:   repo,
		Logger: logger,
		ID:     id,
	}
}

func (w *NotificationWorker) ProcessNextBatch(ctx context.Context) {
	jobs, err := w.Repo.FetchPendingJobs(ctx, 5)
	if err != nil {
		w.Logger.Error().Err(err).Msg("Worker failed to fetch jobs")
		return
	}

	if len(jobs) == 0 {
		return
	}

	w.Logger.Info().Int("count", len(jobs)).Msg("Worker found pending jobs")

	for _, job := range jobs {
		w.processJob(ctx, job)
	}
}

func (w *NotificationWorker) processJob(ctx context.Context, job *repository.Job) {
	w.Logger.Info().Interface("job_id", job.ID).Msg("Processing job...")

	err := w.sendNotification(job)
	if err != nil {
		w.Logger.Warn().Err(err).Interface("job_id", job.ID).Msg("Notification failed, attempting retry logic")

		// This will increment retries and move status to FAILED or PERMANENTLY_FAILED
		retryErr := w.Repo.FailJobWithRetry(ctx, job.ID, 3)
		if retryErr != nil {
			w.Logger.Error().Err(retryErr).Msg("Critical: Could not update failure status in database")
		}
		return
	}

	err = w.Repo.UpdateJobStatus(ctx, job.ID, "SUCCESS")
	if err != nil {
		w.Logger.Error().Err(err).Interface("job_id", job.ID).Msg("Failed to update job status to SUCCESS")
		return
	}

	w.Logger.Info().Interface("job_id", job.ID).Msg("Successfully processed notification")
}

// Helper to simulate real work that might fail
func (w *NotificationWorker) sendNotification(job *repository.Job) error {
	// Simulate processing time
	time.Sleep(500 * time.Millisecond)

	var payload map[string]interface{}
	importJsonErr := json.Unmarshal(job.Payload, &payload)

	if importJsonErr == nil && payload["title"] == "FAIL" {
		return fmt.Errorf("simulated provider downtime")
	}

	return nil
}
