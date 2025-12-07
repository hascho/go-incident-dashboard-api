package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hascho/go-incident-dashboard-api/internal/db"
	"github.com/hascho/go-incident-dashboard-api/internal/handler"
	"github.com/hascho/go-incident-dashboard-api/internal/middleware"
	"github.com/hascho/go-incident-dashboard-api/internal/repository"
	"github.com/hascho/go-incident-dashboard-api/internal/service"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	dbURL := "postgres://user:password@localhost:5432/incidentdb?sslmode=disable"
	dbConfig := db.Config{URL: dbURL}
	dbConn, err := db.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbConn.Close() // ensure the connection is closed when main exits

	logger.Info().Msg("Database connection pool established successfully")

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.LoggerMiddleware(logger))

	incidentRepo := repository.NewIncidentRepository(dbConn)
	incidentService := service.NewIncidentService(incidentRepo)
	incidentHandler := handler.NewIncidentHandler(incidentService)

	r.GET("/incidents/:incidentID", incidentHandler.GetIncidentByID)

	r.POST("/incidents", incidentHandler.CreateIncident)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API is operational"})
	})

	logger.Info().Str("port", "8080").Msg("Server starting")

	r.Run(":8080")
}
