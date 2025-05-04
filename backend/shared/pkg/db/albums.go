package db

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// AlbumStore defines operations specific to albums.
type AlbumStore interface {
	CreateAlbum(ctx context.Context, userID models.UserID, name, description string) (*models.Album, error)
	ListAlbumsByUserID(ctx context.Context, userID models.UserID) ([]models.Album, error)
	GetAlbumByID(ctx context.Context, userID models.UserID, albumID models.AlbumID) (*models.Album, error)         // Changed from models.AlbumID
	UpdateAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID, name, description string) error // Changed from models.AlbumID
	DeleteAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID) error                           // Changed from models.AlbumID

	// Album-Image relationship operations
	AddImageToAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID, imageID models.ImageID) error      // Changed from models.AlbumID
	RemoveImageFromAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID, imageID models.ImageID) error // Changed from models.AlbumID
	ListImagesInAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID) ([]models.ImageMetadata, error)  // Changed from models.AlbumID
}

// --- AlbumStore Implementation ---

// CreateAlbum creates a new album for the specified user
func (s *PostgresStore) CreateAlbum(ctx context.Context, userID models.UserID, name, description string) (*models.Album, error) {
	log.Printf("DB: CreateAlbum called for UserID: %s, Name: %s", userID, name)

	newAlbumID := uuid.New().String()

	query := `
		INSERT INTO albums (id, user_id, name, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	var album models.Album
	err := s.Pool.QueryRow(ctx, query, newAlbumID, userID, name, description).Scan(
		&album.ID, &album.UserID, &album.Name, &album.Description,
		&album.CreatedAt, &album.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating album: %v", err)
		return nil, err
	}

	log.Printf("DB: Successfully created album ID: %s for user ID: %s", album.ID, userID)
	return &album, nil
}

// ListAlbumsByUserID retrieves all albums for a specific user
func (s *PostgresStore) ListAlbumsByUserID(ctx context.Context, userID models.UserID) ([]models.Album, error) {
	log.Printf("DB: ListAlbumsByUserID called for UserID: %s", userID)

	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM albums
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying albums for user %s: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var album models.Album
		err := rows.Scan(
			&album.ID, &album.UserID, &album.Name, &album.Description,
			&album.CreatedAt, &album.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning album row: %v", err)
			return nil, err
		}
		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating album rows for user %s: %v", userID, err)
		return nil, err
	}

	log.Printf("DB: Found %d albums for user ID: %s", len(albums), userID)
	// Return empty slice if no albums found, not nil
	if albums == nil {
		albums = []models.Album{}
	}
	return albums, nil
}

// GetAlbumByID retrieves a specific album by ID, ensuring it belongs to the specified user
func (s *PostgresStore) GetAlbumByID(ctx context.Context, userID models.UserID, albumID models.AlbumID) (*models.Album, error) {
	log.Printf("DB: GetAlbumByID called for UserID: %s, AlbumID: %s", userID, albumID)

	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM albums
		WHERE user_id = $1 AND id = $2
	`

	var album models.Album
	err := s.Pool.QueryRow(ctx, query, userID, albumID).Scan(
		&album.ID, &album.UserID, &album.Name, &album.Description,
		&album.CreatedAt, &album.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("album not found")
		}
		log.Printf("Error getting album: %v", err)
		return nil, err
	}

	return &album, nil
}

// UpdateAlbum updates an existing album's details
func (s *PostgresStore) UpdateAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID, name, description string) error {
	log.Printf("DB: UpdateAlbum called for UserID: %s, AlbumID: %s", userID, albumID)

	query := `
		UPDATE albums
		SET name = $3, description = $4, updated_at = NOW()
		WHERE user_id = $1 AND id = $2
	`

	result, err := s.Pool.Exec(ctx, query, userID, albumID, name, description)
	if err != nil {
		log.Printf("Error updating album: %v", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("album not found")
	}

	log.Printf("DB: Successfully updated album ID: %s", albumID)
	return nil
}

// DeleteAlbum deletes an album and all its image associations
func (s *PostgresStore) DeleteAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID) error {
	log.Printf("DB: DeleteAlbum called for UserID: %s, AlbumID: %s", userID, albumID)

	// Start a transaction to ensure atomicity
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx) // Rollback if not committed

	// First, verify the album exists and belongs to the user
	verifyQuery := `
		SELECT id FROM albums
		WHERE user_id = $1 AND id = $2
	`
	var id models.AlbumID
	err = tx.QueryRow(ctx, verifyQuery, userID, albumID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("album not found")
		}
		log.Printf("Error verifying album ownership: %v", err)
		return err
	}

	// Delete album-image relationships first (this will be handled by CASCADE, but being explicit)
	_, err = tx.Exec(ctx, `DELETE FROM album_images WHERE album_id = $1`, albumID)
	if err != nil {
		log.Printf("Error deleting album image relations: %v", err)
		return err
	}

	// Now delete the album itself
	_, err = tx.Exec(ctx, `DELETE FROM albums WHERE id = $1`, albumID)
	if err != nil {
		log.Printf("Error deleting album: %v", err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("DB: Successfully deleted album ID: %s", albumID)
	return nil
}

// AddImageToAlbum adds an image to an album
func (s *PostgresStore) AddImageToAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID, imageID models.ImageID) error {
	log.Printf("DB: AddImageToAlbum called for UserID: %s, AlbumID: %s, ImageID: %s", userID, albumID, imageID)

	// Start a transaction
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Verify the album exists and belongs to the user
	verifyAlbumQuery := `
		SELECT id FROM albums
		WHERE user_id = $1 AND id = $2
	`
	var albID models.AlbumID
	err = tx.QueryRow(ctx, verifyAlbumQuery, userID, albumID).Scan(&albID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("album not found")
		}
		log.Printf("Error verifying album ownership: %v", err)
		return err
	}

	// Verify the image exists and belongs to the user
	verifyImageQuery := `
		SELECT id FROM images
		WHERE user_id = $1 AND id = $2
	`
	var imgID models.ImageID
	err = tx.QueryRow(ctx, verifyImageQuery, userID, imageID).Scan(&imgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("image not found")
		}
		log.Printf("Error verifying image ownership: %v", err)
		return err
	}

	// Add the image to the album
	insertQuery := `
		INSERT INTO album_images (album_id, image_id)
		VALUES ($1, $2)
		ON CONFLICT (album_id, image_id) DO NOTHING
	`
	_, err = tx.Exec(ctx, insertQuery, albumID, imageID)
	if err != nil {
		log.Printf("Error adding image to album: %v", err)
		return err
	}

	// Update the album's updated_at timestamp
	_, err = tx.Exec(ctx, `UPDATE albums SET updated_at = NOW() WHERE id = $1`, albumID)
	if err != nil {
		log.Printf("Error updating album timestamp: %v", err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("DB: Successfully added image %s to album %s", imageID, albumID)
	return nil
}

// RemoveImageFromAlbum removes an image from an album
func (s *PostgresStore) RemoveImageFromAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID, imageID models.ImageID) error {
	log.Printf("DB: RemoveImageFromAlbum called for UserID: %s, AlbumID: %s, ImageID: %s", userID, albumID, imageID)

	// Start a transaction
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Verify the album exists and belongs to the user
	verifyAlbumQuery := `
		SELECT id FROM albums
		WHERE user_id = $1 AND id = $2
	`
	var albID models.AlbumID
	err = tx.QueryRow(ctx, verifyAlbumQuery, userID, albumID).Scan(&albID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("album not found")
		}
		log.Printf("Error verifying album ownership: %v", err)
		return err
	}

	// Remove the image from the album
	deleteQuery := `
		DELETE FROM album_images
		WHERE album_id = $1 AND image_id = $2
	`
	result, err := tx.Exec(ctx, deleteQuery, albumID, imageID)
	if err != nil {
		log.Printf("Error removing image from album: %v", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("image not found in album")
	}

	// Update the album's updated_at timestamp
	_, err = tx.Exec(ctx, `UPDATE albums SET updated_at = NOW() WHERE id = $1`, albumID)
	if err != nil {
		log.Printf("Error updating album timestamp: %v", err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("DB: Successfully removed image %s from album %s", imageID, albumID)
	return nil
}

// ListImagesInAlbum retrieves all images in a specific album
func (s *PostgresStore) ListImagesInAlbum(ctx context.Context, userID models.UserID, albumID models.AlbumID) ([]models.ImageMetadata, error) {
	log.Printf("DB: ListImagesInAlbum called for UserID: %s, AlbumID: %s", userID, albumID)

	// First verify the album belongs to the user
	verifyQuery := `
		SELECT id FROM albums
		WHERE user_id = $1 AND id = $2
	`
	var albID models.AlbumID
	err := s.Pool.QueryRow(ctx, verifyQuery, userID, albumID).Scan(&albID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("album not found")
		}
		log.Printf("Error verifying album ownership: %v", err)
		return nil, err
	}

	// Then get all images in the album
	query := `
		SELECT i.id, i.user_id, i.filename, i.storage_path, i.content_type,
		       i.size, i.width, i.height, i.created_at, i.updated_at
		FROM images i
		JOIN album_images ai ON i.id = ai.image_id
		WHERE ai.album_id = $1 AND i.user_id = $2
		ORDER BY ai.added_at DESC
	`

	rows, err := s.Pool.Query(ctx, query, albumID, userID)
	if err != nil {
		log.Printf("Error querying images in album %s: %v", albumID, err)
		return nil, err
	}
	defer rows.Close()

	var images []models.ImageMetadata
	for rows.Next() {
		var img models.ImageMetadata
		err := rows.Scan(
			&img.ID, &img.UserID, &img.Filename, &img.StoragePath, &img.ContentType,
			&img.Size, &img.Width, &img.Height, &img.CreatedAt, &img.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning image row: %v", err)
			return nil, err
		}
		images = append(images, img)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating image rows: %v", err)
		return nil, err
	}

	log.Printf("DB: Found %d images in album ID: %s", len(images), albumID)
	// Return empty slice if no images found, not nil
	if images == nil {
		images = []models.ImageMetadata{}
	}
	return images, nil
}
