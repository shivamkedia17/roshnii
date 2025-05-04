package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // For generating unique image IDs/paths

	"io"

	"github.com/shivamkedia17/roshnii/services/server/internal/middleware" // Adjust import paths
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
	"github.com/shivamkedia17/roshnii/shared/pkg/storage" // Add this import
)

// Requires Access to Blob Storage
type ImageHandler struct {
	Config  *config.Config
	DB      db.ImageStore
	Storage storage.BlobStorage
}

func NewImageHandler(config *config.Config, db db.ImageStore, storage storage.BlobStorage) *ImageHandler {
	return &ImageHandler{
		Config:  config,
		DB:      db,
		Storage: storage,
	}
}

// Update HandleUploadImage to store the actual file
func (h *ImageHandler) HandleUploadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Retrieve the file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'file' field in form data"})
		return
	}

	// Basic validation (add more robust checks: size limit, allowed types)
	if fileHeader.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is empty"})
		return
	}

	// Limit file size (e.g., 20MB)
	const maxUploadSize = 20 * 1024 * 1024
	if fileHeader.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File size exceeds limit of %d MB", maxUploadSize/1024/1024)})
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported file type: %s", contentType)})
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded file"})
		return
	}
	defer file.Close()

	// Generate a unique ID for the image
	imageID := uuid.New().String()

	// Upload the file to storage
	storagePath, err := h.Storage.Upload(c.Request.Context(), fileHeader.Filename, userID, file, contentType)
	if err != nil {
		log.Printf("Error uploading file to storage: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store uploaded file"})
		return
	}

	// Create metadata in database
	metadata := &models.ImageMetadata{
		ID:          imageID,
		UserID:      userID,
		Filename:    fileHeader.Filename,
		StoragePath: storagePath,
		ContentType: contentType,
		Size:        fileHeader.Size,
		Width:       0, // TODO: Extract dimensions if needed
		Height:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.DB.CreateImageMetadata(c.Request.Context(), metadata)
	if err != nil {
		// If DB storage fails, try to clean up the file we just uploaded
		cleanupErr := h.Storage.Delete(c.Request.Context(), storagePath)
		if cleanupErr != nil {
			log.Printf("Warning: Failed to clean up file after DB error: %v", cleanupErr)
		}

		log.Printf("Error saving image metadata to DB for image %s: %v", imageID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record image information"})
		return
	}

	log.Printf("Successfully uploaded and saved metadata for image ID: %s", imageID)
	c.JSON(http.StatusCreated, metadata)
}

// Implement HandleDownloadImage to serve the actual file
func (h *ImageHandler) HandleDownloadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Get the image ID from the URL
	imageID := c.Param("id")

	// Get the image metadata from the database
	meta, err := h.DB.GetImageByID(c.Request.Context(), userID, imageID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image metadata"})
		return
	}

	// Get the file from storage
	file, contentType, err := h.Storage.Download(c.Request.Context(), meta.StoragePath)
	if err != nil {
		log.Printf("Error retrieving file from storage: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image file"})
		return
	}
	defer file.Close()

	// Set the response headers
	c.Header("Content-Disposition", "inline; filename="+meta.Filename)
	c.Header("Content-Type", contentType)

	// Stream the file to the response
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, file); err != nil {
		log.Printf("Error streaming file to response: %v", err)
		// Can't really do anything here as we've already started writing the response
	}
}

// HandleGetImage retrieves metadata for a single image.
func (h *ImageHandler) HandleGetImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}
	imageID := c.Param("id")
	meta, err := h.DB.GetImageByID(c.Request.Context(), userID, imageID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image"})
		return
	}
	c.JSON(http.StatusOK, meta)
}

func (h *ImageHandler) HandleDeleteImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Get the image ID from the URL
	imageID := c.Param("id")

	// Get the image metadata from the database
	meta, err := h.DB.GetImageByID(c.Request.Context(), userID, imageID)
	if err != nil {
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image metadata"})
		return
	}

	// Delete the file from storage
	if err := h.Storage.Delete(c.Request.Context(), meta.StoragePath); err != nil {
		log.Printf("Error deleting file from storage: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
		return
	}

	// Delete the metadata from the database
	// We need to implement this method in the PostgresStore
	if err := h.DB.DeleteImageByID(c.Request.Context(), userID, imageID); err != nil {
		log.Printf("Error deleting image metadata from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

// HandleListImages retrieves images for the logged-in user.
func (h *ImageHandler) HandleListImages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Use the actual DB store method
	images, err := h.DB.ListImagesByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error listing images for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve images"})
		return
	}

	// If no images found, the DB function should return an empty slice `[]models.ImageMetadata{}`
	// So no need for explicit nil check here if DB function is correct.

	log.Printf("Retrieved %d images for user %s", len(images), userID)
	c.JSON(http.StatusOK, images)
}
