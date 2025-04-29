package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
)

// TagHandler handles tag-related API requests.
type TagHandler struct {
	Store db.Store
	AppConfig *config.Config
}

func NewTagHandler(store db.Store, cfg *config.Config) *TagHandler {
	return &TagHandler{Store: store, AppConfig: cfg}
}

// RegisterRoutes connects tag routes to the Gin engine.
func (h *TagHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	tagRoutes := router.Group("/tags")
	tagRoutes.Use(authMiddleware)
	{
		tagRoutes.POST("", h.CreateTag)
		tagRoutes.GET("", h.ListTags)
		tagRoutes.POST(":id/images/:image_id", h.AddTagToImage)
		tagRoutes.DELETE(":id/images/:image_id", h.RemoveTagFromImage)
		// Optionally: tagRoutes.DELETE(":id", h.DeleteTag)
	}
}

func (h *TagHandler) CreateTag(c *gin.Context)        {}
func (h *TagHandler) ListTags(c *gin.Context)         {}
func (h *TagHandler) AddTagToImage(c *gin.Context)    {}
func (h *TagHandler) RemoveTagFromImage(c *gin.Context) {}
