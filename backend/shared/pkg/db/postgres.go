package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines the database operations.
type Store interface {
	UserStore
	ImageStore
	AlbumStore
	Close()
}

// PostgresStore holds the connection pool for PostgreSQL interactions.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// NewPostgresStore creates a new PostgreSQL connection pool and returns a Store interface.
func NewPostgresStore(databaseURL string) (Store, error) {
	log.Println("Initializing PostgreSQL connection pool...")

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Error parsing Postgres config: %v", err)
		return nil, err
	}

	// Increase default connections slightly for a web app
	config.MaxConns = 15
	config.MinConns = 3
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("Error creating Postgres connection pool: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		log.Printf("Error pinging Postgres database: %v", err)
		return nil, err
	}

	// Initialize DB schema from file
	err = initializeSchema(ctx, pool)
	if err != nil {
		pool.Close()
		log.Printf("Error initializing database schema: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL and verified connection.")
	return &PostgresStore{Pool: pool}, nil
}

// initializeSchema reads and executes the schema.sql file
func initializeSchema(ctx context.Context, pool *pgxpool.Pool) error {
	// Path relative to where the application runs
	schemaPath := "./schema.sql"

	// Read the schema file
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Printf("Error reading schema file: %v", err)
		// Try a fallback location
		alternatePath := "../db/schema.sql"
		schemaSQL, err = os.ReadFile(alternatePath)
		if err != nil {
			return fmt.Errorf("could not read schema file from %s or %s: %w",
				schemaPath, alternatePath, err)
		}
	}

	// Execute the schema SQL
	_, err = pool.Exec(ctx, string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema SQL: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// Close closes the PostgreSQL connection pool.
func (s *PostgresStore) Close() {
	if s.Pool != nil {
		log.Println("Closing PostgreSQL connection pool...")
		s.Pool.Close()
	}
}

// Helper function to check for common constraint violations
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

func isNotNullViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23502"
}
