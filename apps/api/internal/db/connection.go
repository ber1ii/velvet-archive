package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitPool creates and validates a thread-safe connection pool for Postgres
func InitPool(databaseURL string) (*pgxpool.Pool, error) {
	// Parse the connection configuration string
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Tweak connection pool settings for sanity
	config.MaxConns = 25
	config.MinConns = 3
	config.MaxConnIdleTime = 30 * time.Minute

	// Context with a timeout prevents the server hanging indefinitely if DB is dead
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to ensure the credentials and connection are physically working
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Successfully established thread-safe PostgreSQL connection pool.")
	return pool, nil
}
