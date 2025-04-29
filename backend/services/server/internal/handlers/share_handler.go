package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
)

// ShareHandler handles sharing-related API requests.
type ShareHandler struct {
	Store     db.Store
	AppConfig *config.Config
}

func NewShareHandler(store db.Store, cfg *config.Config) *ShareHandler {
	return &ShareHandler{Store: store, AppConfig: cfg}
}

// RegisterRoutes connects share routes to the Gin engine.
func (h *ShareHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	shareRoutes := router.Group("/shares")
	shareRoutes.Use(authMiddleware)
	{
		shareRoutes.POST("", h.ShareResource)
		shareRoutes.DELETE(":id", h.UnshareResource)
		shareRoutes.GET("", h.ListShares)
		shareRoutes.GET("received", h.ListReceivedShares)
	}
}

func (h *ShareHandler) ShareResource(c *gin.Context)      {}
func (h *ShareHandler) UnshareResource(c *gin.Context)    {}
func (h *ShareHandler) ListShares(c *gin.Context)         {}
func (h *ShareHandler) ListReceivedShares(c *gin.Context) {}
