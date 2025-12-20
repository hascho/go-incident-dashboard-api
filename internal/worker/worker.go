package worker

import (
	"context"
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

	// --- SIMULATING WORK ---
	// In a real app, this is where we would call SendGrid, Twilio, or Slack.
	time.Sleep(1 * time.Second)
	// -----------------------

	err := w.Repo.UpdateJobStatus(ctx, job.ID, "SUCCESS")
	if err != nil {
		w.Logger.Error().Err(err).Interface("job_id", job.ID).Msg("Failed to update job status to SUCCESS")
		return
	}

	w.Logger.Info().Interface("job_id", job.ID).Msg("Successfully processed notification")
}
