package service

import (
	"context"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/repository"
)

type IncidentService interface {
	CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error)
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
