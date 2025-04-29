package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// Store defines the database operations.
type Store interface {
	UserStore
	ImageStore
	Close() // Add Close method to the interface
}

// UserStore defines operations specific to users.
type UserStore interface {
	FindOrCreateUserByGoogleID(ctx context.Context, googleUser *models.GoogleUser) (*models.User, error)
	// Added for dev login:
	FindOrCreateUserByEmail(ctx context.Context, email string, name string, provider string) (*models.User, error)
	// Add other user methods if needed, e.g., GetUserByID
}

// ImageStore defines operations specific to images.
type ImageStore interface {
	CreateImageMetadata(ctx context.Context, meta *models.ImageMetadata) error
	ListImagesByUserID(ctx context.Context, userID models.UserID) ([]models.ImageMetadata, error)
	// GetImageByID retrieves metadata for a given image ID and user.
	GetImageByID(ctx context.Context, userID models.UserID, imageID models.ImageID) (*models.ImageMetadata, error)
}

// PostgresStore holds the connection pool for PostgreSQL interactions.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// NewPostgresStore creates a new PostgreSQL connection pool and returns a Store interface.
func NewPostgresStore(databaseURL string) (Store, error) { // Return Store interface
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

// Helper function to check for common constraint violations
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

func isNotNullViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23502"
}

// --- ImageStore Implementation ---

// CreateImageMetadata inserts metadata about a newly uploaded image.
func (s *PostgresStore) CreateImageMetadata(ctx context.Context, meta *models.ImageMetadata) error {
	log.Printf("DB: CreateImageMetadata called for UserID: %d, Filename: %s, ImageID: %s", meta.UserID, meta.Filename, meta.ID)

	query := `
        INSERT INTO images (id, user_id, filename, storage_path, content_type, size, width, height)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.Pool.Exec(ctx, query,
		meta.ID, meta.UserID, meta.Filename, meta.StoragePath, meta.ContentType,
		meta.Size, meta.Width, meta.Height, // Width/Height can be null if not provided
	)
	if err != nil {
		log.Printf("Error inserting image metadata: %v", err)
		return err
	}
	log.Printf("DB: Successfully inserted metadata for image ID: %s", meta.ID)
	return nil
}

// ListImagesByUserID retrieves all image metadata for a specific user.
func (s *PostgresStore) ListImagesByUserID(ctx context.Context, userID models.UserID) ([]models.ImageMetadata, error) {
	log.Printf("DB: ListImagesByUserID called for UserID: %d", userID)

	query := `
        SELECT id, user_id, filename, storage_path, content_type, size, width, height, created_at, updated_at
        FROM images
        WHERE user_id = $1
        ORDER BY created_at DESC`

	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying images for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var images []models.ImageMetadata
	for rows.Next() {
		var img models.ImageMetadata
		// Nullable columns (width, height, storage_path) need care, though pgx handles *sql.Null types well.
		// Ensure Scan parameters match SELECT order.
		err := rows.Scan(
			&img.ID, &img.UserID, &img.Filename, &img.StoragePath, &img.ContentType,
			&img.Size, &img.Width, &img.Height, &img.CreatedAt, &img.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning image row: %v", err)
			// Decide whether to return partial results or fail completely
			return nil, err
		}
		images = append(images, img)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating image rows for user %d: %v", userID, err)
		return nil, err
	}

	log.Printf("DB: Found %d images for user ID: %d", len(images), userID)
	// Return empty slice if no images found, not nil
	if images == nil {
		images = []models.ImageMetadata{}
	}
	return images, nil
}

// GetImageByID retrieves metadata for a single image belonging to a user.
func (s *PostgresStore) GetImageByID(ctx context.Context, userID models.UserID, imageID models.ImageID) (*models.ImageMetadata, error) {
	log.Printf("DB: GetImageByID called for UserID: %d, ImageID: %s", userID, imageID)

	query := `
        SELECT id, user_id, filename, storage_path, content_type, size, width, height, created_at, updated_at
        FROM images
        WHERE user_id = $1 AND id = $2`

	var img models.ImageMetadata
	err := s.Pool.QueryRow(ctx, query, userID, imageID).Scan(
		&img.ID, &img.UserID, &img.Filename, &img.StoragePath, &img.ContentType,
		&img.Size, &img.Width, &img.Height, &img.CreatedAt, &img.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("image not found")
		}
		return nil, err
	}
	return &img, nil
}
