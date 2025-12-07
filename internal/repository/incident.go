package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hascho/go-incident-dashboard-api/internal/model"
)

type IncidentRepository interface {
	CreateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error)
	GetIncidentByID(ctx context.Context, id string) (*model.Incident, error)
	GetAllIncidents(ctx context.Context) ([]*model.Incident, error)
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

func (r *incidentRepository) GetIncidentByID(ctx context.Context, id string) (*model.Incident, error) {
	query := `
		SELECT id, title, description, status, severity, team, created_at, updated_at
		FROM incidents
		WHERE id = $1`

	incident := &model.Incident{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Description,
		&incident.Status,
		&incident.Severity,
		&incident.Team,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("repository: failed to get incident %s: %w", id, err)
	}

	return incident, nil
}

func (r *incidentRepository) GetAllIncidents(ctx context.Context) ([]*model.Incident, error) {
	query := `
		SELECT id, title, description, status, severity, team, created_at, updated_at
		FROM incidents
		ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query incidents: %w", err)
	}
	defer rows.Close()

	incidents := make([]*model.Incident, 0)

	for rows.Next() {
		incident := &model.Incident{}
		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.Status,
			&incident.Severity,
			&incident.Team,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan incident row: %w", err)
		}
		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows iteration error: %w", err)
	}

	return incidents, nil
}
