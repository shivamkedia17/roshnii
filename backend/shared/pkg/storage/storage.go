package storage

import (
	"context"
	"io"
	"time"

	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// BlobStorage defines the interface for file storage operations
type BlobStorage interface {
	// Upload stores a file and returns its storage path
	Upload(ctx context.Context, filename string, userId models.UserID, content io.Reader, contentType string) (string, error)

	// Download retrieves a file by its storage path
	Download(ctx context.Context, storagePath string) (io.ReadCloser, string, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, storagePath string) error

	// GenerateURL creates a temporary URL for accessing the file (optional)
	GenerateURL(ctx context.Context, storagePath string, expiry time.Duration) (string, error)
}
