package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type contextKey string

const LoggerContextKey contextKey = "logger"

// LoggerMiddleware creates a request-scoped logger and injects it into context
func LoggerMiddleware(baseLogger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get(RequestIDKey)

		requestLogger := baseLogger.With().
			Str("request_id", requestID.(string)).
			Logger()

		ctx := context.WithValue(c.Request.Context(), LoggerContextKey, requestLogger)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetLogger is a helper function to safely retrieve the request-scoped logger
func GetLogger(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(zerolog.Logger); ok {
		return logger
	}
	return zerolog.Nop()
}
