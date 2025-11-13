package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/g3offrey/idiomapi/internal/model"
	"github.com/g3offrey/idiomapi/internal/repository"
	"github.com/g3offrey/idiomapi/internal/service"
	"github.com/gin-gonic/gin"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service *service.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

// CreateTodo handles POST /api/v1/todos
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req model.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	todo, err := h.service.CreateTodo(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create todo",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodo handles GET /api/v1/todos/:id
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid todo ID",
		})
		return
	}

	todo, err := h.service.GetTodo(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get todo",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// ListTodos handles GET /api/v1/todos
func (h *TodoHandler) ListTodos(c *gin.Context) {
	page := 1
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.DefaultQuery("page_size", "10"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = ps
		}
	}

	var completed *bool
	if completedStr := c.Query("completed"); completedStr != "" {
		completedVal := completedStr == "true"
		completed = &completedVal
	}

	response, err := h.service.ListTodos(c.Request.Context(), page, pageSize, completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list todos",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTodo handles PUT /api/v1/todos/:id
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid todo ID",
		})
		return
	}

	var req model.UpdateTodoRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: bindErr.Error(),
		})
		return
	}

	todo, err := h.service.UpdateTodo(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update todo",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo handles DELETE /api/v1/todos/:id
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid todo ID",
		})
		return
	}

	err = h.service.DeleteTodo(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete todo",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
