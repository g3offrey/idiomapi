package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestTodoResponseJSON(t *testing.T) {
	response := TodoResponse{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded TodoResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, response.ID, decoded.ID)
	assert.Equal(t, response.Title, decoded.Title)
	assert.Equal(t, response.Description, decoded.Description)
	assert.Equal(t, response.Completed, decoded.Completed)
}

func TestTodoListResponseJSON(t *testing.T) {
	response := TodoListResponse{
		Todos: []TodoResponse{
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

func TestErrorResponseJSON(t *testing.T) {
	response := ErrorResponse{
		Error:   "validation_error",
		Message: "Field is required",
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded ErrorResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, response.Error, decoded.Error)
	assert.Equal(t, response.Message, decoded.Message)
}
