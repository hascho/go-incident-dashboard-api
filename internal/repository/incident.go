package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
)

type IncidentRepository interface {
	CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error)
}

type incidentRepository struct {
	DB *sql.DB
}

func NewIncidentRepository(db *sql.DB) IncidentRepository {
	return &incidentRepository{DB: db}
}

func (r *incidentRepository) CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error) {
	query := `
		INSERT INTO incidents (title, description, status, severity, team)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	err := r.DB.QueryRowContext(
		ctx,
		query,
		incident.Title,
		incident.Description,
		incident.Status,
		incident.Severity,
		incident.Team,
	).Scan(&incident.ID, &incident.CreatedAt, &incident.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to create incident: %w", err)
	}
	// return the incident struct, now populated with the database-generated fields (ID, dates)
	return incident, nil
}
