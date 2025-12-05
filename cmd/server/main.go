package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hascho/go-incident-dashboard-api/internal/handler"
	"github.com/hascho/go-incident-dashboard-api/internal/middleware"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.LoggerMiddleware(logger))

	incidentHandler := handler.NewIncidentHandler()

	r.GET("/incidents/:incidentID", incidentHandler.GetIncidentByID)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API is operational"})
	})

	logger.Info().Str("port", "8080").Msg("Server starting")

	r.Run(":8080")
}
