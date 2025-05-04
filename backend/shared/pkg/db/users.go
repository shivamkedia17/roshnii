package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

type UserStore interface {
	FindOrCreateUserByGoogleID(ctx context.Context, googleUser *models.GoogleUser) (*models.User, error)
	GetUserByID(ctx context.Context, userID models.UserID) (*models.User, error)

	// Added for dev login:
	FindOrCreateUserByEmail(ctx context.Context, email string, name string, provider string) (*models.User, error)
}

func (s *PostgresStore) FindOrCreateUserByGoogleID(ctx context.Context, googleUser *models.GoogleUser) (*models.User, error) {
	log.Printf("DB: FindOrCreateUserByGoogleID called for Google ID: %s, Email: %s", googleUser.ID, googleUser.Email)

	// First try to find the user by Google ID
	findQuery := `
        SELECT id, google_id, email, name, picture_url, auth_provider, created_at, updated_at
        FROM users WHERE google_id = $1`

	var user models.User
	err := s.Pool.QueryRow(ctx, findQuery, googleUser.ID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == nil {
		// User exists, return them
		return &user, nil
	}

	if err != pgx.ErrNoRows {
		// Unexpected error
		return nil, err
	}

	// User not found, create a new one with a generated UUID
	newID := uuid.New()

	insertQuery := `
        INSERT INTO users (id, google_id, email, name, picture_url, auth_provider)
        VALUES ($1, $2, $3, $4, $5, 'google')
        RETURNING id, google_id, email, name, picture_url, auth_provider, created_at, updated_at`

	err = s.Pool.QueryRow(ctx, insertQuery,
		newID, googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture,
	).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating new user with Google ID: %v", err)
		return nil, err
	}

	log.Printf("DB: Created new user ID: %s for Google ID: %s", user.ID, googleUser.ID)
	return &user, nil
}

// FindOrCreateUserByEmail finds a user by email or creates one (for dev login).
func (s *PostgresStore) FindOrCreateUserByEmail(ctx context.Context, email string, name string, provider string) (*models.User, error) {
	log.Printf("DB: FindOrCreateUserByEmail called for Email: %s", email)

	// Try selecting the user first to see if they exist
	selectQuery := `
        SELECT id, google_id, email, name, picture_url, auth_provider, created_at, updated_at
        FROM users WHERE email = $1`

	var user models.User
	err := s.Pool.QueryRow(ctx, selectQuery, email).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == nil {
		// User exists, return them
		return &user, nil
	}

	if err != pgx.ErrNoRows {
		// Unexpected error
		return nil, err
	}

	// User not found, create a new one with a generated UUID
	newID := uuid.New()

	insertQuery := `
        INSERT INTO users (id, email, name, auth_provider)
        VALUES ($1, $2, $3, $4)
        RETURNING id, google_id, email, name, picture_url, auth_provider, created_at, updated_at`

	err = s.Pool.QueryRow(ctx, insertQuery, newID, email, name, provider).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating new user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (s *PostgresStore) GetUserByID(ctx context.Context, userID models.UserID) (*models.User, error) {
	log.Printf("DB: GetUserByID called for ID: %s", userID)

	query := `
		SELECT id, google_id, email, name, picture_url, auth_provider, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	var googleID sql.NullString

	err := s.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID, &googleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error querying user by ID: %v", err)
		return nil, err
	}

	if googleID.Valid {
		user.GoogleID = &googleID.String
	}

	return &user, nil
}
