package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hascho/go-incident-dashboard-api/internal/middleware"
)

type IncidentHandler struct{}

func NewIncidentHandler() *IncidentHandler {
	return &IncidentHandler{}
}

func (h *IncidentHandler) GetIncidentByID(c *gin.Context) {
	incidentID := c.Param("incidentID")

	logger := middleware.GetLogger(c.Request.Context())

	logger.Info().
		Str("component", "handler").
		Str("incident_id", incidentID).
		Msg("Context-aware log test")

	c.JSON(http.StatusOK, gin.H{
		"incident_id": incidentID,
		"source":      "internal/handler/incident.go",
		"messsage":    "Structured logging injected successfully.",
	})
}
