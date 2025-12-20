package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hascho/go-incident-dashboard-api/internal/db"
	"github.com/hascho/go-incident-dashboard-api/internal/repository"
	"github.com/hascho/go-incident-dashboard-api/internal/worker"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("service", "notification-worker").Logger()

	dbURL := "postgres://user:password@127.0.0.1:5432/incidentdb?sslmode=disable"
	dbConfig := db.Config{URL: dbURL}
	dbConn, err := db.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Worker failed to connect to database")
	}
	defer dbConn.Close()

	jobRepo := repository.NewJobRepository(dbConn)
	notificationWorker := worker.NewNotificationWorker(jobRepo, logger, "worker-01")

	// The Heartbeat (Polling Loop)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msg("Worker Service Started. Checking for jobs...")

	// This is the infinite loop that keeps our worker alive
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Wake up and work!
				notificationWorker.ProcessNextBatch(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for someone to kill the process (Ctrl+C)
	<-quit
	logger.Warn().Msg("Worker received shutdown signal. Stopping...")
	cancel()

	// Give the worker a moment to finish its current job
	time.Sleep(2 * time.Second)
	logger.Info().Msg("Worker exited.")
}
