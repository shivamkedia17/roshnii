package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
)

// SearchHandler handles search-related API requests.
type SearchHandler struct {
	Store     db.Store
	AppConfig *config.Config
}

func NewSearchHandler(store db.Store, cfg *config.Config) *SearchHandler {
	return &SearchHandler{Store: store, AppConfig: cfg}
}

func (h *SearchHandler) Search(c *gin.Context) {}
