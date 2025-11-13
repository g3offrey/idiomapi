package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTodoModel(t *testing.T) {
	now := time.Now()
	todo := Todo{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, 1, todo.ID)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.Equal(t, "Test Description", todo.Description)
	assert.False(t, todo.Completed)
	assert.Equal(t, now, todo.CreatedAt)
	assert.Equal(t, now, todo.UpdatedAt)
}
