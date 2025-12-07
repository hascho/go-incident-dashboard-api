package service

import (
	"context"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/repository"
)

type IncidentService interface {
	CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error)
	GetIncidentByID(ctx context.Context, id string) (*model.Incident, error)
	GetAllIncidents(ctx context.Context) ([]*model.Incident, error)
}

type incidentService struct {
	Repo repository.IncidentRepository
}

func NewIncidentService(repo repository.IncidentRepository) IncidentService {
	return &incidentService{Repo: repo}
}

func (s *incidentService) CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error) {
	// todo add business logic here
	return s.Repo.CreateIncident(ctx, incident)
}

func (s *incidentService) GetIncidentByID(ctx context.Context, id string) (*model.Incident, error) {
	return s.Repo.GetIncidentByID(ctx, id)
}

func (s *incidentService) GetAllIncidents(ctx context.Context) ([]*model.Incident, error) {
	return s.Repo.GetAllIncidents(ctx)
}
