package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodoJSON(t *testing.T) {
	todo := Todo{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
	}

	// Marshal to JSON
	data, err := json.Marshal(todo)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Unmarshal from JSON
	var decoded Todo
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, todo.ID, decoded.ID)
	assert.Equal(t, todo.Title, decoded.Title)
	assert.Equal(t, todo.Description, decoded.Description)
	assert.Equal(t, todo.Completed, decoded.Completed)
}

func TestCreateTodoRequestJSON(t *testing.T) {
	req := CreateTodoRequest{
		Title:       "New Todo",
		Description: "New Description",
		Completed:   false,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded CreateTodoRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, req.Title, decoded.Title)
	assert.Equal(t, req.Description, decoded.Description)
	assert.Equal(t, req.Completed, decoded.Completed)
}

func TestUpdateTodoRequestJSON(t *testing.T) {
	title := "Updated Title"
	completed := true

	req := UpdateTodoRequest{
		Title:     &title,
		Completed: &completed,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded UpdateTodoRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.NotNil(t, decoded.Title)
	assert.Equal(t, title, *decoded.Title)
	assert.NotNil(t, decoded.Completed)
	assert.Equal(t, completed, *decoded.Completed)
}

func TestTodoListResponseJSON(t *testing.T) {
	response := TodoListResponse{
		Todos: []Todo{
			{ID: 1, Title: "Todo 1", Completed: false},
			{ID: 2, Title: "Todo 2", Completed: true},
		},
		Total:      2,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded TodoListResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Len(t, decoded.Todos, 2)
	assert.Equal(t, response.Total, decoded.Total)
	assert.Equal(t, response.Page, decoded.Page)
}
