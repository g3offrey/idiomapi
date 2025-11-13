package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/g3offrey/idiomapi/internal/dto"
	"github.com/g3offrey/idiomapi/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrNotFound is returned when a todo is not found
	ErrNotFound = errors.New("todo not found")
)

// TodoRepository handles todo data operations
type TodoRepository struct {
	pool *pgxpool.Pool
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(pool *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{pool: pool}
}

// Create creates a new todo
func (r *TodoRepository) Create(ctx context.Context, req dto.CreateTodoRequest) (*model.Todo, error) {
	query := `
		INSERT INTO todos (title, description, completed)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, completed, created_at, updated_at
	`

	var todo model.Todo
	err := r.pool.QueryRow(ctx, query, req.Title, req.Description, req.Completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return &todo, nil
}

// GetByID retrieves a todo by its ID
func (r *TodoRepository) GetByID(ctx context.Context, id int) (*model.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo model.Todo
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// List retrieves a paginated list of todos
func (r *TodoRepository) List(ctx context.Context, page, pageSize int, completed *bool) ([]model.Todo, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Build query based on filters
	var countQuery, listQuery string
	var args []interface{}

	if completed != nil {
		countQuery = "SELECT COUNT(*) FROM todos WHERE completed = $1"
		listQuery = `
			SELECT id, title, description, completed, created_at, updated_at
			FROM todos
			WHERE completed = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = append(args, *completed, pageSize, offset)
	} else {
		countQuery = "SELECT COUNT(*) FROM todos"
		listQuery = `
			SELECT id, title, description, completed, created_at, updated_at
			FROM todos
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = append(args, pageSize, offset)
	}

	// Get total count
	var total int
	if completed != nil {
		err := r.pool.QueryRow(ctx, countQuery, *completed).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count todos: %w", err)
		}
	} else {
		err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count todos: %w", err)
		}
	}

	// Get todos
	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, total, nil
}

// Update updates a todo
func (r *TodoRepository) Update(ctx context.Context, id int, req dto.UpdateTodoRequest) (*model.Todo, error) {
	// First check if todo exists
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	query := "UPDATE todos SET "
	args := []interface{}{}
	argPosition := 1
	updates := []string{}

	if req.Title != nil {
		updates = append(updates, fmt.Sprintf("title = $%d", argPosition))
		args = append(args, *req.Title)
		argPosition++
	}

	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPosition))
		args = append(args, *req.Description)
		argPosition++
	}

	if req.Completed != nil {
		updates = append(updates, fmt.Sprintf("completed = $%d", argPosition))
		args = append(args, *req.Completed)
		argPosition++
	}

	if len(updates) == 0 {
		// No fields to update, return existing
		return existing, nil
	}

	query += fmt.Sprintf("%s WHERE id = $%d RETURNING id, title, description, completed, created_at, updated_at",
		joinStrings(updates, ", "), argPosition)
	args = append(args, id)

	var todo model.Todo
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return &todo, nil
}

// Delete deletes a todo by ID
func (r *TodoRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM todos WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
