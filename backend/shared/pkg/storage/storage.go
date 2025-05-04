package storage

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/shivamkedia17/roshnii/shared/pkg/config"
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

type StorageType string

const (
	Local = "local"
	// S3,
	// etc.
)

func InitStorage(cfg *config.Config) (BlobStorage, error) {
	var storageService BlobStorage
	var err error

	storeType := cfg.BlobStorageType

	if storeType == Local {
		localStoragePath := cfg.LocalstoragePath
		storageService, err = NewLocalStorage(localStoragePath)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		log.Printf("Using local file storage at: %s", localStoragePath)

		return storageService, nil
	} else {
		// For now, fall back to local storage if type is unrecognized
		log.Printf("Unrecognized storage type '%s', using local storage", cfg.BlobStorageType)
		storageService, _ = NewLocalStorage("./uploads")
	}

	return nil, err
}
