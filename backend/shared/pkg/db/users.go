package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// UserStore defines operations specific to users.
type UserStore interface {
	FindOrCreateUserByGoogleID(ctx context.Context, googleUser *models.GoogleUser) (*models.User, error)
	// Added for dev login:
	FindOrCreateUserByEmail(ctx context.Context, email string, name string, provider string) (*models.User, error)
	GetUserByID(ctx context.Context, userID models.UserID) (*models.User, error) // Add this line
}

// --- UserStore Implementation ---

// FindOrCreateUserByGoogleID finds a user by their Google ID or creates a new one.
func (s *PostgresStore) FindOrCreateUserByGoogleID(ctx context.Context, googleUser *models.GoogleUser) (*models.User, error) {
	log.Printf("DB: FindOrCreateUserByGoogleID called for Google ID: %s, Email: %s", googleUser.ID, googleUser.Email)

	var user models.User
	// Use INSERT ... ON CONFLICT to handle find or create atomically
	query := `
        INSERT INTO users (google_id, email, name, picture_url, auth_provider)
        VALUES ($1, $2, $3, $4, 'google')
        ON CONFLICT (google_id) DO UPDATE
        SET email = EXCLUDED.email, name = EXCLUDED.name, picture_url = EXCLUDED.picture_url, updated_at = NOW()
        WHERE users.google_id = $1 -- Ensure the conflict target matches the WHERE clause
        RETURNING id, google_id, email, name, picture_url, auth_provider, created_at, updated_at`

	err := s.Pool.QueryRow(ctx, query,
		googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture,
	).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error in FindOrCreateUserByGoogleID query: %v", err)
		return nil, err
	}
	log.Printf("DB: Found or created user ID: %d for Google ID: %s", user.ID, googleUser.ID)
	return &user, nil
}

// FindOrCreateUserByEmail finds a user by email or creates one (for dev login).
func (s *PostgresStore) FindOrCreateUserByEmail(ctx context.Context, email string, name string, provider string) (*models.User, error) {
	log.Printf("DB: FindOrCreateUserByEmail called for Email: %s", email)

	var user models.User
	// Use INSERT ... ON CONFLICT on the email unique constraint
	query := `
        INSERT INTO users (email, name, auth_provider)
        VALUES ($1, $2, $3)
        ON CONFLICT (email) DO UPDATE
        SET name = EXCLUDED.name, -- Optionally update name on conflict
            updated_at = NOW()
        WHERE users.email = $1 -- Ensure the conflict target matches the WHERE clause
        RETURNING id, google_id, email, name, picture_url, auth_provider, created_at, updated_at`

	// Note: google_id might be null here if created via this method
	err := s.Pool.QueryRow(ctx, query, email, name, provider).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
		&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		// Handle case where user exists but was created via Google (has google_id)
		// and we tried inserting without it. A simple SELECT might be better in this case.
		if errors.Is(err, pgx.ErrNoRows) || isForeignKeyViolation(err) || isNotNullViolation(err) {
			// Try selecting the user first if insert fails ambiguously
			log.Printf("DB: Insert conflict/error for %s, attempting select.", email)
			selectQuery := `
                SELECT id, google_id, email, name, picture_url, auth_provider, created_at, updated_at
                FROM users WHERE email = $1`
			err = s.Pool.QueryRow(ctx, selectQuery, email).Scan(
				&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL,
				&user.AuthProvider, &user.CreatedAt, &user.UpdatedAt,
			)
			if err != nil {
				log.Printf("Error in FindOrCreateUserByEmail (select fallback): %v", err)
				return nil, err
			}
		} else {
			log.Printf("Error in FindOrCreateUserByEmail query: %v", err)
			return nil, err
		}
	}
	log.Printf("DB: Found or created user ID: %d for Email: %s", user.ID, email)
	return &user, nil
}

func (s *PostgresStore) GetUserByID(ctx context.Context, userID models.UserID) (*models.User, error) {
	log.Printf("DB: GetUserByID called for ID: %d", userID)

	query := `
		SELECT id, google_id, email, name, picture_url, auth_provider, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	var googleID sql.NullString // For NULL google_id

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
