package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

const defaultPath = "./uploads"

// LocalStorage implements BlobStorage using the local filesystem
type LocalStorage struct {
	BasePath string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if basePath == "" {
		basePath = defaultPath
	}

	// Ensure the base directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{BasePath: basePath}, nil
}

// Upload stores a file on the local filesystem
func (s *LocalStorage) Upload(ctx context.Context, filename string, userId models.UserID, content io.Reader, contentType string) (string, error) {
	// Create a user-specific directory
	userDir := filepath.Join(s.BasePath, fmt.Sprintf("user_%s", userId))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	// Generate a unique filename (you can improve this with UUID)
	timestamp := time.Now().UnixNano()
	uniqueFilename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(filename))
	storagePath := filepath.Join(fmt.Sprintf("user_%s", userId), uniqueFilename)
	fullPath := filepath.Join(s.BasePath, storagePath)

	// Create the file
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Copy the content to the file
	if _, err := io.Copy(f, content); err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return storagePath, nil
}

// Download retrieves a file from the local filesystem
func (s *LocalStorage) Download(ctx context.Context, storagePath string) (io.ReadCloser, string, error) {
	fullPath := filepath.Join(s.BasePath, storagePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("file not found: %s", storagePath)
	}

	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}

	// Get the content type (simplistic approach - in production use a more robust method)
	contentType := "application/octet-stream" // default
	ext := filepath.Ext(fullPath)
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	return file, contentType, nil
}

// Delete removes a file from the local filesystem
func (s *LocalStorage) Delete(ctx context.Context, storagePath string) error {
	fullPath := filepath.Join(s.BasePath, storagePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // Already deleted, not an error
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GenerateURL creates a URL for direct file access (for local storage, just return the path)
func (s *LocalStorage) GenerateURL(ctx context.Context, storagePath string, expiry time.Duration) (string, error) {
	// For local storage, we'll return a local URL path
	return "/api/image/" + filepath.Base(storagePath) + "/download", nil
}
