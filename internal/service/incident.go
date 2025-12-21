package service

import (
	"context"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/queue"
	"github.com/hascho/go-incident-dashboard-api/internal/repository"
	"github.com/rs/zerolog"
)

type IncidentService interface {
	CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error)
	GetIncidentByID(ctx context.Context, id string) (*model.Incident, error)
	GetAllIncidents(ctx context.Context) ([]*model.Incident, error)
	UpdateIncident(ctx context.Context, incidentID string, req model.UpdateIncidentRequest) (*model.Incident, error)
	DeleteIncident(ctx context.Context, incidentID string) error
}

type incidentService struct {
	Repo    repository.IncidentRepository
	Logger  zerolog.Logger
	JobRepo repository.JobRepository
	Queue   queue.TaskQueue
}

func NewIncidentService(repo repository.IncidentRepository, logger zerolog.Logger, jobRepo repository.JobRepository, q queue.TaskQueue) IncidentService {
	return &incidentService{
		Repo:    repo,
		Logger:  logger,
		JobRepo: jobRepo,
		Queue:   q,
	}
}

func (s *incidentService) CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error) {
	createdIncident, err := s.Repo.CreateIncident(ctx, incident)
	if err != nil {
		return nil, err
	}

	if err := s.JobRepo.CreateJob(ctx, createdIncident); err != nil {
		// Log the failure but do NOT fail the HTTP request
		s.Logger.Error().Err(err).Msg("Failed to create job entry for notification. Incident created but notification is missing.")
	}

	if err := s.Queue.Publish(ctx, createdIncident.ID); err != nil {
		// if Redis is down, it's okay! polling will catch it.
		s.Logger.Warn().Err(err).Msg("Redis publish failed - worker will catch up via polling")
	} else {
		s.Logger.Info().Str("incident_id", createdIncident.ID).Msg("Published job to Redis")
	}

	return createdIncident, nil
}

func (s *incidentService) GetIncidentByID(ctx context.Context, id string) (*model.Incident, error) {
	return s.Repo.GetIncidentByID(ctx, id)
}

func (s *incidentService) GetAllIncidents(ctx context.Context) ([]*model.Incident, error) {
	return s.Repo.GetAllIncidents(ctx)
}

func (s *incidentService) UpdateIncident(ctx context.Context, incidentID string, req model.UpdateIncidentRequest) (*model.Incident, error) {
	existingIncident, err := s.Repo.GetIncidentByID(ctx, incidentID)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		existingIncident.Status = *req.Status
	}
	if req.Description != nil {
		existingIncident.Description = *req.Description
	}

	return s.Repo.UpdateIncident(ctx, existingIncident)
}

func (s *incidentService) DeleteIncident(ctx context.Context, incidentID string) error {
	return s.Repo.DeleteIncident(ctx, incidentID)
}
