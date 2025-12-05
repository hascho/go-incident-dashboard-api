package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	URL string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	conn, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	conn.SetMaxOpenConns(25)                 // max connections open at once
	conn.SetMaxIdleConns(25)                 // max connections waiting idle
	conn.SetConnMaxLifetime(5 * time.Minute) // how long a connection can be reused
	conn.SetConnMaxIdleTime(1 * time.Minute) // how long an idle connection can wait

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}
