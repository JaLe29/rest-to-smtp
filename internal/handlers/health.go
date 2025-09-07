package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rest-to-smtp/internal/models"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health provides a simple health check endpoint
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:  "healthy",
		Service: "rest-to-smtp",
	})
}
