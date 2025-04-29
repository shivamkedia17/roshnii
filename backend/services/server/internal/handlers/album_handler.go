package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
)

// AlbumHandler handles album-related API requests.
type AlbumHandler struct {
	Store db.Store
	AppConfig *config.Config
}

func NewAlbumHandler(store db.Store, cfg *config.Config) *AlbumHandler {
	return &AlbumHandler{Store: store, AppConfig: cfg}
}

// RegisterRoutes connects album routes to the Gin engine.
func (h *AlbumHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	albumRoutes := router.Group("/albums")
	albumRoutes.Use(authMiddleware)
	{
		albumRoutes.POST("", h.CreateAlbum)
		albumRoutes.GET("", h.ListAlbums)
		albumRoutes.GET(":id", h.GetAlbum)
		albumRoutes.DELETE(":id", h.DeleteAlbum)
		albumRoutes.POST(":id/images", h.AddImageToAlbum)
		albumRoutes.DELETE(":id/images/:image_id", h.RemoveImageFromAlbum)
	}
}

func (h *AlbumHandler) CreateAlbum(c *gin.Context)    {}
func (h *AlbumHandler) ListAlbums(c *gin.Context)     {}
func (h *AlbumHandler) GetAlbum(c *gin.Context)       {}
func (h *AlbumHandler) DeleteAlbum(c *gin.Context)    {}
func (h *AlbumHandler) AddImageToAlbum(c *gin.Context) {}
func (h *AlbumHandler) RemoveImageFromAlbum(c *gin.Context) {}
