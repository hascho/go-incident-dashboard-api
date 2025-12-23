package model

import "time"

type Incident struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	Severity           string    `json:"severity"`
	Team               string    `json:"team"`
	NotificationStatus string    `json:"notification_status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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
