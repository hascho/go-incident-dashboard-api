package service

import (
	"context"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/repository"
	"github.com/hascho/go-incident-dashboard-api/internal/worker"
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
	Repo         repository.IncidentRepository
	Logger       zerolog.Logger
	JobProcessor *worker.JobProcessor
}

func NewIncidentService(repo repository.IncidentRepository, logger zerolog.Logger, jp *worker.JobProcessor) IncidentService {
	return &incidentService{Repo: repo, Logger: logger, JobProcessor: jp}
}

func (s *incidentService) CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error) {
	createdIncident, err := s.Repo.CreateIncident(ctx, incident)
	if err != nil {
		return nil, err
	}

	job := worker.Job{
		Incident: createdIncident,
		Logger:   s.Logger,
	}

	// send the job to the buffered channel.
	// if the channel buffer (1000 slots) is full, this line will BLOCK the HTTP request
	// until a worker finishes a job. This is the Backpressure mechanism.
	s.JobProcessor.JobQueue <- job

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
