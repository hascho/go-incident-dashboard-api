package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	dbURL := "postgres://user:password@127.0.0.1:5432/incidentdb?sslmode=disable"
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
	jobRepo := repository.NewJobRepository(dbConn)

	incidentService := service.NewIncidentService(incidentRepo, logger, jobRepo)
	incidentHandler := handler.NewIncidentHandler(incidentService)

	r.GET("/incidents", incidentHandler.GetAllIncidents)
	r.GET("/incidents/:id", incidentHandler.GetIncidentByID)
	r.POST("/incidents", incidentHandler.CreateIncident)
	r.PATCH("/incidents/:id", incidentHandler.PatchIncident)
	r.DELETE("/incidents/:id", incidentHandler.DeleteIncident)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API is operational"})
	})

	logger.Info().Str("port", "8080").Msg("Server starting")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// start the server in a goroutine so main() can proceed to the shutdown block
	go func() {
		logger.Info().Str("port", "8080").Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server forced to close")
		}
	}()

	// GRACEFUL SHUTDOWN BLOCK
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // block execution until a signal is received
	logger.Warn().Msg("Server received shutdown signal. Initiating graceful shutdown...")

	// stop the HTTP server from accepting new requests
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer httpCancel()
	if err := srv.Shutdown(httpCtx); err != nil {
		logger.Fatal().Err(err).Msg("HTTP Server forced to shutdown")
	}

	logger.Info().Msg("Server exiting.")
}
