package dto

import (
	"testing"
	"time"

	"github.com/g3offrey/idiomapi/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestToTodoResponse(t *testing.T) {
	now := time.Now()
	todo := &model.Todo{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := ToTodoResponse(todo)

	assert.Equal(t, todo.ID, response.ID)
	assert.Equal(t, todo.Title, response.Title)
	assert.Equal(t, todo.Description, response.Description)
	assert.Equal(t, todo.Completed, response.Completed)
	assert.Equal(t, todo.CreatedAt, response.CreatedAt)
	assert.Equal(t, todo.UpdatedAt, response.UpdatedAt)
}

func TestToTodoResponseList(t *testing.T) {
	now := time.Now()
	todos := []model.Todo{
		{
			ID:          1,
			Title:       "Todo 1",
			Description: "Description 1",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Title:       "Todo 2",
			Description: "Description 2",
			Completed:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	responses := ToTodoResponseList(todos)

	assert.Len(t, responses, 2)
	assert.Equal(t, todos[0].ID, responses[0].ID)
	assert.Equal(t, todos[0].Title, responses[0].Title)
	assert.Equal(t, todos[1].ID, responses[1].ID)
	assert.Equal(t, todos[1].Title, responses[1].Title)
}

func TestToTodoListResponse(t *testing.T) {
	now := time.Now()
	todos := []model.Todo{
		{
			ID:          1,
			Title:       "Todo 1",
			Description: "Description 1",
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	response := ToTodoListResponse(todos, 10, 1, 5)

	assert.Len(t, response.Todos, 1)
	assert.Equal(t, 10, response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 5, response.PageSize)
	assert.Equal(t, 2, response.TotalPages) // 10 items / 5 per page = 2 pages
}

func TestToTodoListResponse_EmptyList(t *testing.T) {
	todos := []model.Todo{}

	response := ToTodoListResponse(todos, 0, 1, 10)

	assert.Len(t, response.Todos, 0)
	assert.Equal(t, 0, response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, 1, response.TotalPages) // Minimum 1 page
}
