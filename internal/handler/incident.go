package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hascho/go-incident-dashboard-api/internal/middleware"
	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/service"
)

type IncidentHandler struct {
	Service service.IncidentService
}

func NewIncidentHandler(svc service.IncidentService) *IncidentHandler {
	return &IncidentHandler{Service: svc}
}

func (h *IncidentHandler) GetIncidentByID(c *gin.Context) {
	logger := middleware.GetLogger(c.Request.Context())
	incidentID := c.Param("id")

	if _, err := uuid.Parse(incidentID); err != nil {
		logger.Warn().Str("incident_id", incidentID).Msg("Received malformed incident ID")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Incident ID '%s' is not a valid UUID format.", incidentID),
		})
		return
	}

	incident, err := h.Service.GetIncidentByID(c.Request.Context(), incidentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Str("incident_id", incidentID).Msg("Incident not found")
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Incident with ID %s not found", incidentID),
			})
			return
		}

		logger.Error().Err(err).Msg("Failed to retrieve incident")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	response := model.IncidentResponse{
		ID:       incidentID,
		Title:    incident.ID,
		Status:   incident.Status,
		Severity: incident.Severity,
		Team:     incident.Team,
	}

	c.JSON(http.StatusOK, response)
}

func (h *IncidentHandler) CreateIncident(c *gin.Context) {
	logger := middleware.GetLogger(c.Request.Context())

	var req model.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid request: %v", err.Error()),
		})
		return
	}

	incident := &model.Incident{
		Title:       req.Title,
		Description: req.Description,
		Severity:    req.Severity,
		Team:        req.Team,
		Status:      "open", // Handler dictates the initial state
	}

	createdIncident, err := h.Service.CreateIncident(c.Request.Context(), incident)
	if err != nil {
		logger.Error().Err(err).Msg("Repository failure during incident creation")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create incident due to database error.",
		})
		return
	}

	response := model.IncidentResponse{
		ID:       createdIncident.ID,
		Title:    createdIncident.Title,
		Status:   createdIncident.Status,
		Severity: createdIncident.Severity,
		Team:     createdIncident.Team,
	}

	logger.Info().Str("incident_id", createdIncident.ID).Msg("Incident created successfully")
	c.JSON(http.StatusCreated, response)
}

func (h *IncidentHandler) GetAllIncidents(c *gin.Context) {
	logger := middleware.GetLogger(c.Request.Context())

	incidents, err := h.Service.GetAllIncidents(c.Request.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list incidents")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve incident list",
		})
		return
	}

	response := make([]model.IncidentResponse, len(incidents))
	for i, incident := range incidents {
		response[i] = model.IncidentResponse{
			ID:       incident.ID,
			Title:    incident.Title,
			Status:   incident.Status,
			Severity: incident.Severity,
			Team:     incident.Team,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *IncidentHandler) PatchIncident(c *gin.Context) {
	logger := middleware.GetLogger(c.Request.Context())
	incidentID := c.Param("id")

	if _, err := uuid.Parse(incidentID); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Incident ID '%s' is not a valid UUID format.", incidentID),
		})
		return
	}

	var req model.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid request payload: %v", err.Error()),
		})
		return
	}

	updatedIncident, err := h.Service.UpdateIncident(c.Request.Context(), incidentID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Incident with ID %s not found", incidentID),
			})
			return
		}
		logger.Error().Err(err).Msg("Failed to update incident")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error during update.",
		})
		return
	}

	response := model.IncidentResponse{
		ID:       updatedIncident.ID,
		Title:    updatedIncident.Title,
		Status:   updatedIncident.Status,
		Severity: updatedIncident.Severity,
		Team:     updatedIncident.Team,
	}

	c.JSON(http.StatusOK, response)
}
