package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/services/server/internal/auth"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAlbumStore is a mock implementation of the AlbumStore interface for testing
type MockAlbumStore struct {
	albums      map[int64]models.Album
	albumImages map[int64][]string // Album ID -> slice of image IDs
}

func NewMockAlbumStore() *MockAlbumStore {
	return &MockAlbumStore{
		albums:      make(map[int64]models.Album),
		albumImages: make(map[int64][]string),
	}
}

// Implement AlbumStore methods
func (m *MockAlbumStore) CreateAlbum(ctx context.Context, userID models.UserID, name, description string) (*models.Album, error) {
	// Generate a unique ID for the new album
	newID := int64(len(m.albums) + 1)

	// Create album with current time
	album := models.Album{
		ID:          newID,
		UserID:      userID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Store in our mock database
	m.albums[newID] = album

	// Initialize empty image list for this album
	m.albumImages[newID] = []string{}

	return &album, nil
}

func (m *MockAlbumStore) ListAlbumsByUserID(ctx context.Context, userID models.UserID) ([]models.Album, error) {
	var userAlbums []models.Album

	for _, album := range m.albums {
		if album.UserID == userID {
			userAlbums = append(userAlbums, album)
		}
	}

	return userAlbums, nil
}

func (m *MockAlbumStore) GetAlbumByID(ctx context.Context, userID models.UserID, albumID int64) (*models.Album, error) {
	album, exists := m.albums[albumID]
	if !exists {
		return nil, errors.New("album not found")
	}

	// Verify the album belongs to the user
	if album.UserID != userID {
		return nil, errors.New("album not found")
	}

	return &album, nil
}

func (m *MockAlbumStore) UpdateAlbum(ctx context.Context, userID models.UserID, albumID int64, name, description string) error {
	album, exists := m.albums[albumID]
	if !exists {
		return errors.New("album not found")
	}

	// Verify the album belongs to the user
	if album.UserID != userID {
		return errors.New("album not found")
	}

	// Update album
	album.Name = name
	album.Description = description
	album.UpdatedAt = time.Now()

	// Store updated album
	m.albums[albumID] = album

	return nil
}

func (m *MockAlbumStore) DeleteAlbum(ctx context.Context, userID models.UserID, albumID int64) error {
	album, exists := m.albums[albumID]
	if !exists {
		return errors.New("album not found")
	}

	// Verify the album belongs to the user
	if album.UserID != userID {
		return errors.New("album not found")
	}

	// Delete album and its image references
	delete(m.albums, albumID)
	delete(m.albumImages, albumID)

	return nil
}

func (m *MockAlbumStore) AddImageToAlbum(ctx context.Context, userID models.UserID, albumID int64, imageID models.ImageID) error {
	album, exists := m.albums[albumID]
	if !exists {
		return errors.New("album not found")
	}

	// Verify the album belongs to the user
	if album.UserID != userID {
		return errors.New("album not found")
	}

	// Add image to album
	images := m.albumImages[albumID]
	for _, id := range images {
		if id == imageID {
			// Image already in album
			return nil
		}
	}

	m.albumImages[albumID] = append(images, imageID)

	// Update album timestamp
	album.UpdatedAt = time.Now()
	m.albums[albumID] = album

	return nil
}

func (m *MockAlbumStore) RemoveImageFromAlbum(ctx context.Context, userID models.UserID, albumID int64, imageID models.ImageID) error {
	album, exists := m.albums[albumID]
	if !exists {
		return errors.New("album not found")
	}

	// Verify the album belongs to the user
	if album.UserID != userID {
		return errors.New("album not found")
	}

	// Find and remove the image
	images := m.albumImages[albumID]
	found := false

	for i, id := range images {
		if id == imageID {
			// Remove this image
			m.albumImages[albumID] = append(images[:i], images[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("image not found in album")
	}

	// Update album timestamp
	album.UpdatedAt = time.Now()
	m.albums[albumID] = album

	return nil
}

func (m *MockAlbumStore) ListImagesInAlbum(ctx context.Context, userID models.UserID, albumID int64) ([]models.ImageMetadata, error) {
	album, exists := m.albums[albumID]
	if !exists {
		return nil, errors.New("album not found")
	}

	// Verify the album belongs to the user
	if album.UserID != userID {
		return nil, errors.New("album not found")
	}

	// In a real implementation, we would query the database for these images
	// Here, we'll return empty metadata since this is just a mock
	var images []models.ImageMetadata

	for _, imageID := range m.albumImages[albumID] {
		images = append(images, models.ImageMetadata{
			ID:        imageID,
			UserID:    userID,
			Filename:  "mock-image.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return images, nil
}

func TestAlbumHandlerCreateAlbum(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockStore := NewMockAlbumStore()
	handler := NewAlbumHandler(mockStore, nil)

	router := gin.Default()
	// Mock authentication middleware
	authMW := func(c *gin.Context) {
		// Set a mock user ID
		c.Set(middleware.UserContextKey, &auth.Claims{UserID: 1, Email: "test@example.com"})
		c.Next()
	}

	api := router.Group("/api")
	albumsGroup := api.Group("/albums")
	albumsGroup.Use(authMW)
	albumsGroup.POST("", handler.CreateAlbum)

	// Test body
	reqBody := `{"name":"Test Album","description":"This is a test album"}`

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/api/albums", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response
	var response models.Album
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response
	assert.NotZero(t, response.ID)
	assert.Equal(t, "Test Album", response.Name)
	assert.Equal(t, "This is a test album", response.Description)
	assert.Equal(t, models.UserID(1), response.UserID)
}
