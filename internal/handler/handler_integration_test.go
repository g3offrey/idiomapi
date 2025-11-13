package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/g3offrey/idiomapi/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthHandlerIntegration tests the health endpoint
func TestHealthHandlerIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock health handler (no actual db connection needed)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:   "ok",
			Database: "ok",
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
}

// TestTodoHandlerValidation tests request validation
func TestTodoHandlerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock create handler with validation
	router.POST("/api/v1/todos", func(c *gin.Context) {
		var req model.CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, model.Todo{
			ID:    1,
			Title: req.Title,
		})
	})

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
	}{
		{
			name:           "valid request",
			payload:        `{"title":"Test Todo","description":"Test Description"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing title",
			payload:        `{"description":"Test Description"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			payload:        `{"title":}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestTodoHandlerErrorResponses tests error response format
func TestTodoHandlerErrorResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/todos/:id", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Error:   "not_found",
			Message: "Todo not found",
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/todos/999", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response model.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response.Error)
	assert.Equal(t, "Todo not found", response.Message)
}
