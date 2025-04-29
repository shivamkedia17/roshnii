package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // For generating unique image IDs/paths

	"github.com/shivamkedia17/roshnii/services/server/internal/middleware" // Adjust import paths
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// ImageHandler handles image-related API requests.
type ImageHandler struct {
	Store db.ImageStore
	// Storage  storage.BlobStorage // Add your blob storage interface here later
	AppConfig *config.Config
}

// NewImageHandler creates a new ImageHandler.
func NewImageHandler(store db.ImageStore /*storage storage.BlobStorage,*/, cfg *config.Config) *ImageHandler {
	return &ImageHandler{
		Store: store,
		// Storage: storage,
		AppConfig: cfg,
	}
}

// RegisterRoutes connects image routes to the Gin engine.
func (h *ImageHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Upload endpoint
	router.POST("/upload", authMiddleware, h.HandleUploadImage)
	// List user images
	router.GET("/images", authMiddleware, h.HandleListImages)
	// Single image metadata
	router.GET("/image/:id", authMiddleware, h.HandleGetImage)
	// Download image file
	router.GET("/image/:id/download", authMiddleware, h.HandleDownloadImage)
	// Delete image
	router.DELETE("/image/:id", authMiddleware, h.HandleDeleteImage)
}

// HandleUploadImage processes image uploads.
func (h *ImageHandler) HandleUploadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		// This should ideally not happen if middleware is working, but check anyway
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Retrieve the file from the form data
	fileHeader, err := c.FormFile("file") // "file" matches the openapi.yaml spec
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

	log.Printf("Received upload: Filename=%s, Size=%d, ContentType=%s, UserID=%d",
		fileHeader.Filename, fileHeader.Size, contentType, userID)

	// Generate a unique ID and storage path for the image
	imageUUID := uuid.New()       // Generate UUID object
	imageID := imageUUID.String() // Get string representation
	fileExt := filepath.Ext(fileHeader.Filename)
	// Example storage path structure: user_<user_id>/<uuid>.<ext>
	// Even if not storing locally now, save this path in metadata.
	storagePath := fmt.Sprintf("user_%d/%s%s", userID, imageID, fileExt)

	// --- TODO: Implement Actual File Storage ---
	// Open the file just to ensure it's readable, but we won't save it for this MVP step.
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening uploaded file header: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded file"})
		return
	}
	defer func(file multipart.File) {
		_ = file.Close() // Best effort close
	}(file)

	// In a real scenario, you would upload `file` to blob storage using `storagePath` here.
	log.Printf("Placeholder: File [%s] would be saved to storage path: %s", fileHeader.Filename, storagePath)

	// --- Store Metadata in Database ---
	// Extract Width/Height (optional, requires image decoding library like "image")
	// e.g., imgConfig, _, err := image.DecodeConfig(file)
	// For now, we'll leave them as 0/null.
	metadata := &models.ImageMetadata{
		ID:          imageID, // Use the generated UUID string
		UserID:      userID,
		Filename:    fileHeader.Filename,
		StoragePath: storagePath,
		ContentType: contentType,
		Size:        fileHeader.Size,
		Width:       0,          // TODO: Populate later if needed
		Height:      0,          // TODO: Populate later if needed
		CreatedAt:   time.Now(), // Let DB handle default
		UpdatedAt:   time.Now(), // Let DB handle default or trigger
	}

	// Use the actual DB store method
	err = h.Store.CreateImageMetadata(c.Request.Context(), metadata)
	if err != nil {
		log.Printf("Error saving image metadata to DB for image %s: %v", imageID, err)
		// TODO: Consider deleting the uploaded file from storage if DB fails (rollback)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record image information"})
		return
	}

	log.Printf("Successfully saved metadata for image ID: %s", imageID)

	// Fetch the created metadata back from DB to include DB-generated timestamps
	// For simplicity in MVP, just return the metadata we constructed.
	// In production, you might query it back or rely on RETURNING clause in SQL.
	c.JSON(http.StatusCreated, metadata) // Return the metadata of the created image record
}

// HandleListImages retrieves images for the logged-in user.
func (h *ImageHandler) HandleListImages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Use the actual DB store method
	images, err := h.Store.ListImagesByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error listing images for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve images"})
		return
	}

	// If no images found, the DB function should return an empty slice `[]models.ImageMetadata{}`
	// So no need for explicit nil check here if DB function is correct.

	log.Printf("Retrieved %d images for user %d", len(images), userID)
	c.JSON(http.StatusOK, images)
}

// HandleGetImage retrieves metadata for a single image.
func (h *ImageHandler) HandleGetImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}
	imageID := c.Param("id")
	meta, err := h.Store.GetImageByID(c.Request.Context(), userID, imageID)
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

func (h *ImageHandler) HandleDeleteImage(c *gin.Context)   {}
func (h *ImageHandler) HandleDownloadImage(c *gin.Context) {}
