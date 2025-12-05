package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hascho/go-incident-dashboard-api/internal/middleware"
	"github.com/hascho/go-incident-dashboard-api/internal/model"
	"github.com/hascho/go-incident-dashboard-api/internal/util"
)

type IncidentHandler struct{}

func NewIncidentHandler() *IncidentHandler {
	return &IncidentHandler{}
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
