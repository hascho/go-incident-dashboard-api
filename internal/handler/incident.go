package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hascho/go-incident-dashboard-api/internal/middleware"
	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/service"
	"github.com/hascho/go-incident-dashboard-api/internal/util"
)

type IncidentHandler struct {
	Service service.IncidentService
}

func NewIncidentHandler(svc service.IncidentService) *IncidentHandler {
	return &IncidentHandler{Service: svc}
}

func (h *IncidentHandler) GetIncidentByID(c *gin.Context) {
	incidentID := c.Param("incidentID")
	logger := middleware.GetLogger(c.Request.Context())

	if incidentID == "not-found" {
		apiErr := util.NewNotFoundError(fmt.Sprintf("Incident with ID %s not found.", incidentID))
		logger.Error().Str("incident_id", incidentID).Msg(apiErr.Error())

		c.JSON(apiErr.StatusCode, model.ErrorResponse{
			Code:    apiErr.StatusCode,
			Message: apiErr.Message,
		})
		return
	}

	logger.Info().
		Str("component", "handler").
		Str("incident_id", incidentID).
		Msg("Context-aware log test")

	response := model.IncidentResponse{
		ID:       incidentID,
		Title:    "DB Error on Prod Cluster",
		Status:   "Acknowledged",
		Severity: "Critical",
		Team:     "SRE",
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
