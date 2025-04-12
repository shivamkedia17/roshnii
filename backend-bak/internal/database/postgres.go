package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore holds the connection pool for PostgreSQL interactions.
type PostgresStore struct {
	// Pool is the pgx connection pool.
	Pool *pgxpool.Pool
}

// NewPostgresStore creates a new PostgreSQL connection pool.
func NewPostgresStore(databaseURL string) (*PostgresStore, error) {
	log.Println("Initializing PostgreSQL connection pool...")

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Error parsing Postgres config: %v", err)
		return nil, err
	}

	// Optional: Configure pool settings here if needed
	// config.MaxConns = 10 // Example: Set max connections
	// config.MinConns = 2  // Example: Set min connections
	// config.MaxConnIdleTime = 5 * time.Minute // Example

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("Error creating Postgres connection pool: %v", err)
		return nil, err
	}

	// Verify the connection with a Ping and timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5-second timeout
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		// If ping fails, close the pool before returning error
		pool.Close()
		log.Printf("Error pinging Postgres database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL and verified connection.")
	return &PostgresStore{Pool: pool}, nil
}

// Close closes the PostgreSQL connection pool.
func (s *PostgresStore) Close() {
	if s.Pool != nil {
		log.Println("Closing PostgreSQL connection pool...")
		s.Pool.Close()
	}
}

// --- Placeholder for Database Interaction Methods ---
// Example (will be implemented later):
// func (s *PostgresStore) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
//     // ... query logic using s.Pool ...
//     return nil, nil
// }
// func (s *PostgresStore) CreateImageMetadata(ctx context.Context, meta *models.ImageMetadata) error {
// 	   // ... insert logic using s.Pool ...
//     return nil
// }
