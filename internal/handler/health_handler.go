package handler

import (
	"net/http"

	"github.com/g3offrey/idiomapi/internal/database"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "ok"
	if err := h.db.Health(c.Request.Context()); err != nil {
		dbStatus = "error"
	}

	status := "ok"
	statusCode := http.StatusOK
	if dbStatus != "ok" {
		status = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, HealthResponse{
		Status:   status,
		Database: dbStatus,
	})
}
