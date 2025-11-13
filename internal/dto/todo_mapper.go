package dto

import "github.com/g3offrey/idiomapi/internal/model"

// ToTodoResponse converts a domain Todo to a TodoResponse DTO
func ToTodoResponse(todo *model.Todo) TodoResponse {
	return TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

// ToTodoResponseList converts a slice of domain Todos to TodoResponse DTOs
func ToTodoResponseList(todos []model.Todo) []TodoResponse {
	responses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = ToTodoResponse(&todo)
	}
	return responses
}

// ToTodoListResponse converts domain data to a TodoListResponse DTO
func ToTodoListResponse(todos []model.Todo, total, page, pageSize int) TodoListResponse {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return TodoListResponse{
		Todos:      ToTodoResponseList(todos),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
