# Contributing to idiomapi

Thank you for your interest in contributing to idiomapi! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful, inclusive, and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker & Docker Compose (optional)
- Make (optional, for using Makefile commands)
- Git

### Setting Up Development Environment

1. **Fork and clone the repository**
```bash
git clone https://github.com/YOUR_USERNAME/idiomapi.git
cd idiomapi
```

2. **Install dependencies**
```bash
go mod download
```

3. **Install development tools**
```bash
make install-tools
```

4. **Install git hooks**
```bash
make install-hooks
```

5. **Start PostgreSQL**
```bash
make docker-up
```

6. **Run migrations**
```bash
make migrate-up
```

7. **Run the application**
```bash
make run
```

## Development Workflow

### 1. Create a Branch

Create a feature branch from `main`:
```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or changes

### 2. Make Your Changes

- Write clean, idiomatic Go code
- Follow existing code style and patterns
- Add tests for new functionality
- Update documentation as needed
- Keep commits focused and atomic

### 3. Test Your Changes

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

Run race detector:
```bash
go test -race ./...
```

### 4. Format and Lint

Format code:
```bash
make fmt
```

Run linter:
```bash
make lint
```

Run go vet:
```bash
make vet
```

Or run all checks at once:
```bash
make ci
```

### 5. Commit Your Changes

Follow conventional commit messages:
```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process or auxiliary tool changes

Example:
```bash
git commit -m "feat(todos): add filtering by completion status"
```

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Create a pull request on GitHub with:
- Clear title and description
- Reference any related issues
- List of changes made
- Screenshots (if applicable)

## Code Style Guidelines

### Go Code Style

1. **Follow Go conventions**
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

2. **Naming conventions**
   - Use camelCase for unexported names
   - Use PascalCase for exported names
   - Use descriptive names
   - Avoid abbreviations unless common

3. **Error handling**
   - Always check errors
   - Wrap errors with context
   - Return errors, don't panic (except in init/main)

4. **Comments**
   - Comment all exported functions/types
   - Use complete sentences
   - Explain "why", not "what"

### Example Code

```go
// CreateTodo creates a new todo item in the database.
// It validates the input and returns an error if validation fails.
func (s *TodoService) CreateTodo(ctx context.Context, req model.CreateTodoRequest) (*model.Todo, error) {
    s.logger.Debug("creating todo", "title", req.Title)
    
    todo, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create todo: %w", err)
    }
    
    s.logger.Info("todo created", "id", todo.ID)
    return todo, nil
}
```

## Testing Guidelines

### Unit Tests

- Test individual functions/methods
- Use table-driven tests when appropriate
- Mock dependencies
- Aim for high coverage

Example:
```go
func TestTodoService_CreateTodo(t *testing.T) {
    tests := []struct {
        name    string
        input   model.CreateTodoRequest
        wantErr bool
    }{
        {
            name: "valid todo",
            input: model.CreateTodoRequest{
                Title: "Test",
            },
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Test HTTP endpoints
- Test database operations
- Use test database or mocks
- Clean up after tests

## Pull Request Process

1. **Ensure CI passes**
   - All tests pass
   - Linting passes
   - No security vulnerabilities

2. **Update documentation**
   - README.md if needed
   - API documentation
   - Code comments

3. **Request review**
   - Tag relevant maintainers
   - Respond to feedback
   - Make requested changes

4. **Merge requirements**
   - At least one approval
   - All checks passing
   - Conflicts resolved

## Project Structure

When adding new code, follow the existing structure:

```
internal/
├── config/      # Configuration management
├── database/    # Database connection
├── handler/     # HTTP handlers
├── middleware/  # HTTP middleware
├── model/       # Data models
├── repository/  # Data access
└── service/     # Business logic
```

## Common Tasks

### Adding a New Endpoint

1. Add model in `internal/model/`
2. Add repository method in `internal/repository/`
3. Add service method in `internal/service/`
4. Add handler in `internal/handler/`
5. Register route in `cmd/api/main.go`
6. Add tests for each layer
7. Update API documentation

### Adding Database Migration

1. Create new file in `migrations/`
2. Name: `XXX_description.sql`
3. Include both up and down migrations
4. Test migration locally
5. Update documentation

### Adding Configuration Option

1. Add field to config struct in `internal/config/config.go`
2. Add to `configs/config.toml`
3. Update documentation
4. Add tests

## Debugging

### Enable Debug Logging

```toml
[logging]
level = "debug"
```

### Run with Delve Debugger

```bash
dlv debug cmd/api/main.go -- -config configs/config.toml
```

### Database Debugging

```bash
psql -h localhost -U postgres -d tododb
```

## Getting Help

- Check existing issues and documentation
- Ask questions in pull request comments
- Create a new issue for bugs or feature requests

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project README

Thank you for contributing to idiomapi! 🎉
