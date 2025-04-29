package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSearchRoutesRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewSearchHandler(nil, nil)
	h.RegisterRoutes(r.Group("/api"), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search", nil)
	r.ServeHTTP(w, req)
	assert.NotEqual(t, 404, w.Code, "Route should be registered")
}
