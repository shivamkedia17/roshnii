package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
)

// AlbumHandler handles album-related API requests.
type AlbumHandler struct {
	Config *config.Config
	DB     db.AlbumStore
}

// NewAlbumHandler creates a new AlbumHandler
func NewAlbumHandler(config *config.Config, db db.AlbumStore) *AlbumHandler {
	return &AlbumHandler{
		Config: config,
		DB:     db,
	}
}

// CreateAlbum handles the creation of a new album
func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Bind request body
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	album, err := h.DB.CreateAlbum(c.Request.Context(), userID, req.Name, req.Description)
	if err != nil {
		log.Printf("Error creating album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	c.JSON(http.StatusCreated, album)
}

// ListAlbums retrieves all albums for the authenticated user
func (h *AlbumHandler) ListAlbums(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	albums, err := h.DB.ListAlbumsByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error listing albums for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums: " + err.Error()})
		return
	}

	log.Printf("Successfully retrieved %d albums for user %s", len(albums), userID)
	c.JSON(http.StatusOK, albums)
}

// GetAlbum retrieves a specific album by ID
func (h *AlbumHandler) GetAlbum(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Parse album ID from URL
	albumID := c.Param("id")
	album, err := h.DB.GetAlbumByID(c.Request.Context(), userID, albumID)
	if err != nil {
		if err.Error() == "album not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		log.Printf("Error getting album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// UpdateAlbum updates an existing album
func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Parse album ID from URL
	albumID := c.Param("id")

	// Bind request body
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	err := h.DB.UpdateAlbum(c.Request.Context(), userID, albumID, req.Name, req.Description)
	if err != nil {
		if err.Error() == "album not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		log.Printf("Error updating album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
}

// DeleteAlbum deletes an album
func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Parse album ID from URL
	albumID := c.Param("id")

	err := h.DB.DeleteAlbum(c.Request.Context(), userID, albumID)
	if err != nil {
		if err.Error() == "album not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		log.Printf("Error deleting album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// AddImageToAlbum adds an image to an album
func (h *AlbumHandler) AddImageToAlbum(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Parse album ID from URL
	albumID := c.Param("id")

	// Bind request body
	var req struct {
		ImageID string `json:"image_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	err := h.DB.AddImageToAlbum(c.Request.Context(), userID, albumID, req.ImageID)
	if err != nil {
		if err.Error() == "album not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		if err.Error() == "image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		log.Printf("Error adding image to album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add image to album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image added to album successfully"})
}

// RemoveImageFromAlbum removes an image from an album
func (h *AlbumHandler) RemoveImageFromAlbum(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Parse album ID and image ID from URL
	albumID := c.Param("id")

	imageID := c.Param("image_id")
	if imageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	err := h.DB.RemoveImageFromAlbum(c.Request.Context(), userID, albumID, imageID)
	if err != nil {
		if err.Error() == "album not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		if err.Error() == "image not found in album" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found in album"})
			return
		}
		log.Printf("Error removing image from album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove image from album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image removed from album successfully"})
}

// ListAlbumImages retrieves all images in a specific album
func (h *AlbumHandler) ListAlbumImages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	albumID := c.Param("id")

	images, err := h.DB.ListImagesInAlbum(c.Request.Context(), userID, albumID)
	if err != nil {
		if err.Error() == "album not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		log.Printf("Error listing images in album: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve images"})
		return
	}

	c.JSON(http.StatusOK, images)
}
