package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// ImageStore defines operations specific to images.
type ImageStore interface {
	CreateImageMetadata(ctx context.Context, meta *models.ImageMetadata) error
	ListImagesByUserID(ctx context.Context, userID models.UserID) ([]models.ImageMetadata, error)
	GetImageByID(ctx context.Context, userID models.UserID, imageID models.ImageID) (*models.ImageMetadata, error)
	DeleteImageByID(ctx context.Context, userID models.UserID, imageID models.ImageID) error // Add this line
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

// DeleteImageByID removes image metadata from the database
func (s *PostgresStore) DeleteImageByID(ctx context.Context, userID models.UserID, imageID models.ImageID) error {
	log.Printf("DB: DeleteImageByID called for UserID: %d, ImageID: %s", userID, imageID)

	query := `
		DELETE FROM images
		WHERE user_id = $1 AND id = $2
	`

	result, err := s.Pool.Exec(ctx, query, userID, imageID)
	if err != nil {
		log.Printf("Error deleting image metadata: %v", err)
		return err
	}

	// Check if any rows were affected
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	log.Printf("DB: Successfully deleted metadata for image ID: %s", imageID)
	return nil
}
