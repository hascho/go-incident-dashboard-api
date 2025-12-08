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
	UpdateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error)
	DeleteIncident(ctx context.Context, id string) error
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

func (r *incidentRepository) UpdateIncident(ctx context.Context, incident *model.Incident) (*model.Incident, error) {
	query := `
		UPDATE incidents
		SET status = $2, description = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, status, severity, team, created_at, updated_at`

	updatedIncident := &model.Incident{}

	err := r.DB.QueryRowContext(ctx, query, incident.ID, incident.Status, incident.Description).Scan(
		&updatedIncident.ID,
		&updatedIncident.Title,
		&updatedIncident.Description,
		&updatedIncident.Status,
		&updatedIncident.Severity,
		&updatedIncident.Team,
		&updatedIncident.CreatedAt,
		&updatedIncident.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("repository: failed to update incident %s: %w", incident.ID, err)
	}

	return updatedIncident, nil
}

func (r *incidentRepository) DeleteIncident(ctx context.Context, id string) error {
	query := `DELETE FROM incidents WHERE id = $1`

	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete incident %s: %w", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to check rows affected: %w", err)
	}

	// rowsAffected is 0 when the ID did not exist
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
