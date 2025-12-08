package model

import "time"

type Incident struct {
	ID          string
	Title       string
	Description string
	Status      string
	Severity    string
	Team        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateIncidentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Severity    string `json:"severity" binding:"required"`
	Team        string `json:"team" binding:"required"`
}

type IncidentResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Team     string `json:"team"`
}

type UpdateIncidentRequest struct {
	Status      *string `json:"status"`
	Description *string `json:"description"`
}
