package service

import (
	"context"
	"log/slog"

	"github.com/g3offrey/idiomapi/internal/dto"
	"github.com/g3offrey/idiomapi/internal/model"
	"github.com/g3offrey/idiomapi/internal/repository"
)

// TodoService handles business logic for todos
type TodoService struct {
	repo   *repository.TodoRepository
	logger *slog.Logger
}

// NewTodoService creates a new TodoService
func NewTodoService(repo *repository.TodoRepository, logger *slog.Logger) *TodoService {
	return &TodoService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(ctx context.Context, req dto.CreateTodoRequest) (*model.Todo, error) {
	s.logger.Debug("creating todo", "title", req.Title)
	todo, err := s.repo.Create(ctx, req)
	if err != nil {
		s.logger.Error("failed to create todo", "error", err)
		return nil, err
	}
	s.logger.Info("todo created", "id", todo.ID, "title", todo.Title)
	return todo, nil
}

// GetTodo retrieves a todo by ID
func (s *TodoService) GetTodo(ctx context.Context, id int) (*model.Todo, error) {
	s.logger.Debug("getting todo", "id", id)
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get todo", "id", id, "error", err)
		return nil, err
	}
	return todo, nil
}

// ListTodos retrieves a paginated list of todos
func (s *TodoService) ListTodos(ctx context.Context, page, pageSize int, completed *bool) ([]model.Todo, int, error) {
	s.logger.Debug("listing todos", "page", page, "pageSize", pageSize)

	todos, total, err := s.repo.List(ctx, page, pageSize, completed)
	if err != nil {
		s.logger.Error("failed to list todos", "error", err)
		return nil, 0, err
	}

	return todos, total, nil
}

// UpdateTodo updates a todo
func (s *TodoService) UpdateTodo(ctx context.Context, id int, req dto.UpdateTodoRequest) (*model.Todo, error) {
	s.logger.Debug("updating todo", "id", id)
	todo, err := s.repo.Update(ctx, id, req)
	if err != nil {
		s.logger.Error("failed to update todo", "id", id, "error", err)
		return nil, err
	}
	s.logger.Info("todo updated", "id", todo.ID)
	return todo, nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	s.logger.Debug("deleting todo", "id", id)
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete todo", "id", id, "error", err)
		return err
	}
	s.logger.Info("todo deleted", "id", id)
	return nil
}
