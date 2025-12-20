package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hascho/go-incident-dashboard-api/internal/model"
)

type Job struct {
	ID         uuid.UUID
	IncidentID uuid.UUID
	Status     string
	Payload    json.RawMessage
	Retries    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type JobRepository interface {
	CreateJob(ctx context.Context, incident *model.Incident) error
	FetchPendingJobs(ctx context.Context, limit int) ([]*Job, error)
	UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status string) error
	FailJobWithRetry(ctx context.Context, jobID uuid.UUID, maxRetries int) error
}

type jobRepository struct {
	DB *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepository{DB: db}
}

func (r *jobRepository) CreateJob(ctx context.Context, incident *model.Incident) error {
	payload, err := json.Marshal(incident)
	if err != nil {
		return fmt.Errorf("failed to marshal incident payload: %w", err)
	}

	query := `
		INSERT INTO notification_jobs (incident_id, payload, status)
		VALUES ($1, $2, 'PENDING')`

	_, err = r.DB.ExecContext(ctx, query, incident.ID, payload)
	if err != nil {
		return fmt.Errorf("repository: failed to insert job: %w", err)
	}
	return nil
}

func (r *jobRepository) FetchPendingJobs(ctx context.Context, limit int) ([]*Job, error) {
	query := `
		SELECT id, incident_id, payload, retries, created_at, updated_at, status
		FROM notification_jobs
		WHERE status = 'PENDING' OR (status = 'FAILED' AND retries < 3)
		ORDER BY created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT $1`

	rows, err := r.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		job := &Job{}
		err := rows.Scan(
			&job.ID,
			&job.IncidentID,
			&job.Payload,
			&job.Retries,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning job row: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *jobRepository) UpdateJobStatus(ctx context.Context, jobID uuid.UUID, status string) error {
	query := `
		UPDATE notification_jobs
		SET status = $2, updated_at = NOW()
		WHERE id = $1`

	_, err := r.DB.ExecContext(ctx, query, jobID, status)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}
	return nil
}

func (r *jobRepository) FailJobWithRetry(ctx context.Context, jobID uuid.UUID, maxRetries int) error {
	query := `
		UPDATE notification_jobs
		SET retries = retries + 1,
			status = CASE
				WHEN retries + 1 >= $2 THEN 'PERMANENTLY_FAILED'
				ELSE 'FAILED'
			END,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.DB.ExecContext(ctx, query, jobID, maxRetries)
	if err != nil {
		return fmt.Errorf("failed to update job retry count: %w", err)
	}
	return nil
}
