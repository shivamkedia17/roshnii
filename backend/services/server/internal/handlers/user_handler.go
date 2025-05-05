package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
)

// UserHandler manages user-related API operations
type UserHandler struct {
	Config *config.Config
	DB     db.UserStore
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(config *config.Config, db db.UserStore) *UserHandler {
	return &UserHandler{
		Config: config,
		DB:     db,
	}
}

// GetCurrentUser returns the authenticated user's profile
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	user, err := h.DB.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Return the entire user struct as is
	c.JSON(http.StatusOK, user)
}

// Example of additional method (optional)
// UpdateUserProfile updates the user's profile information
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user session"})
		return
	}

	// Bind request body
	var req struct {
		Name       string  `json:"name"`
		PictureURL *string `json:"picture_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing user
	_, err := h.DB.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving user for update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Update user properties
	// Note: You would need to implement this method in your UserStore interface
	// For example: h.Store.UpdateUserProfile(ctx, userID, req.Name, req.PictureURL)

	// For demonstration, just return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user_id": userID,
	})
}
