# Architecture Documentation

## Overview

This is an enterprise-ready REST API for a TODO application built following Go best practices and idioms. The application follows clean architecture principles with clear separation of concerns.

## Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go 1.21+ | Modern, performant, statically typed |
| **Web Framework** | Gin | High-performance HTTP web framework |
| **Database Driver** | pgx/v5 | Pure Go PostgreSQL driver with connection pooling |
| **Database** | PostgreSQL 12+ | Reliable, ACID-compliant relational database |
| **Configuration** | cleanenv | TOML/YAML/ENV configuration management |
| **Logging** | slog | Structured logging (Go 1.21+ standard library) |
| **Testing** | testify/assert | Rich assertion library for testing |
| **Linting** | golangci-lint | Fast Go linters aggregator |

## Project Structure

```
idiomapi/
├── cmd/
│   └── api/              # Application entrypoint
│       └── main.go       # Main application initialization
│
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   │
│   ├── database/        # Database connection and setup
│   │   └── database.go
│   │
│   ├── handler/         # HTTP request handlers
│   │   ├── todo_handler.go
│   │   ├── health_handler.go
│   │   └── handler_integration_test.go
│   │
│   ├── middleware/      # HTTP middleware
│   │   ├── logger.go    # Request logging
│   │   └── recovery.go  # Panic recovery
│   │
│   ├── model/           # Data models and DTOs
│   │   ├── todo.go
│   │   └── todo_test.go
│   │
│   ├── repository/      # Data access layer
│   │   ├── todo_repository.go
│   │   └── todo_repository_test.go
│   │
│   └── service/         # Business logic layer
│       └── todo_service.go
│
├── pkg/                 # Public, reusable packages
│   └── logger/          # Logging utilities
│       ├── logger.go
│       └── logger_test.go
│
├── migrations/          # Database migrations
│   └── 001_create_todos_table.sql
│
├── configs/             # Configuration files
│   └── config.toml
│
├── scripts/             # Utility scripts
│   ├── pre-commit.sh    # Git pre-commit hook
│   └── api-example.sh   # API usage examples
│
├── .github/
│   └── workflows/
│       └── ci.yml       # GitHub Actions CI/CD
│
├── Makefile             # Build automation
├── docker-compose.yml   # Local development services
├── .golangci.yml        # Linting configuration
├── .editorconfig        # Editor configuration
├── .gitignore          # Git ignore rules
├── LICENSE             # MIT License
└── README.md           # Project documentation
```

## Architecture Layers

### 1. Handler Layer (`internal/handler/`)

**Responsibility**: Handle HTTP requests and responses

- Parse HTTP requests
- Validate input data
- Call service layer
- Format HTTP responses
- Handle HTTP-specific errors (4xx, 5xx)

**Key Files**:
- `todo_handler.go` - CRUD operations for todos
- `health_handler.go` - Health check endpoint

### 2. Service Layer (`internal/service/`)

**Responsibility**: Business logic and orchestration

- Implement business rules
- Coordinate between repositories
- Transaction management
- Business-level logging
- Business-level error handling

**Key Files**:
- `todo_service.go` - Todo business logic

### 3. Repository Layer (`internal/repository/`)

**Responsibility**: Data access and persistence

- Database queries
- Data mapping
- Handle database-specific errors
- Implement repository pattern

**Key Files**:
- `todo_repository.go` - Todo data access

### 4. Model Layer (`internal/model/`)

**Responsibility**: Data structures and DTOs

- Define domain models
- Request/Response DTOs
- Validation tags
- JSON serialization

**Key Files**:
- `todo.go` - Todo models and DTOs

### 5. Middleware Layer (`internal/middleware/`)

**Responsibility**: Cross-cutting concerns

- Request logging
- Error recovery
- Authentication (if needed)
- CORS handling (if needed)

**Key Files**:
- `logger.go` - Request/response logging
- `recovery.go` - Panic recovery

## Data Flow

```
HTTP Request
    ↓
[Middleware] → Logging, Recovery
    ↓
[Handler] → Parse request, validate
    ↓
[Service] → Business logic
    ↓
[Repository] → Database operations
    ↓
[Database] → PostgreSQL
    ↓
[Repository] → Map results
    ↓
[Service] → Process results
    ↓
[Handler] → Format response
    ↓
[Middleware] → Log response
    ↓
HTTP Response
```

## Design Patterns

### 1. **Clean Architecture**
- Clear separation of concerns
- Dependency injection
- Testable components
- Independent of frameworks

### 2. **Repository Pattern**
- Abstract data access
- Centralize data logic
- Easier to test
- Database-agnostic interface

### 3. **Service Layer Pattern**
- Encapsulate business logic
- Coordinate multiple repositories
- Transaction management
- Reusable business operations

### 4. **DTO Pattern**
- Separate internal models from API contracts
- Version API independently
- Validation at API boundary

## Database Schema

### Todos Table

```sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_todos_completed ON todos(completed);
CREATE INDEX idx_todos_created_at ON todos(created_at);
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/todos` | Create todo |
| GET | `/api/v1/todos` | List todos (with pagination) |
| GET | `/api/v1/todos/:id` | Get todo by ID |
| PUT | `/api/v1/todos/:id` | Update todo |
| DELETE | `/api/v1/todos/:id` | Delete todo |

## Configuration

Configuration is managed through TOML files with support for environment variable overrides.

```toml
[server]
host = "0.0.0.0"
port = 8080

[database]
host = "localhost"
port = 5432
user = "postgres"
password = "postgres"
dbname = "tododb"

[logging]
level = "info"
format = "json"
```

## Logging

Structured logging using Go's standard `slog` package:

- **Levels**: debug, info, warn, error
- **Formats**: JSON, text
- **Context**: Request ID, user info, etc.

Example log entry:
```json
{
  "time": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "msg": "request processed",
  "method": "GET",
  "path": "/api/v1/todos",
  "status": 200,
  "latency": "15ms"
}
```

## Error Handling

Consistent error responses across the API:

```json
{
  "error": "error_code",
  "message": "Human-readable error message"
}
```

## Testing Strategy

### Unit Tests
- Test individual functions/methods
- Mock dependencies
- Fast execution
- High coverage

### Integration Tests
- Test HTTP endpoints
- Mock database when possible
- Test middleware
- Validate request/response

### Example:
```go
func TestCreateTodo(t *testing.T) {
    // Arrange
    handler := NewTodoHandler(mockService)
    
    // Act
    response := handler.CreateTodo(request)
    
    // Assert
    assert.Equal(t, http.StatusCreated, response.Status)
}
```

## Performance Considerations

1. **Connection Pooling**: pgx pool with configurable size
2. **Indexes**: Database indexes on frequently queried fields
3. **Pagination**: Limit result sets to prevent memory issues
4. **Context Propagation**: Timeout and cancellation support
5. **Graceful Shutdown**: Clean connection closure

## Security Features

1. **Input Validation**: Request validation using struct tags
2. **SQL Injection Prevention**: Parameterized queries with pgx
3. **Error Information**: Don't leak internal details
4. **Health Checks**: Monitor application health
5. **Graceful Shutdown**: Proper resource cleanup

## Development Workflow

1. **Write Code**: Follow Go conventions
2. **Run Tests**: `make test`
3. **Format Code**: `make fmt`
4. **Lint**: `make lint`
5. **Build**: `make build`
6. **Run Locally**: `make run`
7. **Pre-commit Hook**: Automatic CI checks

## Deployment

### Local Development
```bash
make docker-up    # Start PostgreSQL
make migrate-up   # Run migrations
make run          # Start application
```

### Production Deployment
```bash
make build        # Build binary
./bin/idiomapi -config /etc/idiomapi/config.toml
```

## Monitoring

- **Health Endpoint**: `/health`
- **Structured Logs**: JSON format for log aggregation
- **Metrics**: Ready for Prometheus integration
- **Database Health**: Connection pool monitoring

## Future Enhancements

1. **Authentication/Authorization**: JWT or OAuth2
2. **Rate Limiting**: Prevent abuse
3. **Caching**: Redis for frequently accessed data
4. **Metrics**: Prometheus/Grafana
5. **Tracing**: OpenTelemetry
6. **API Versioning**: v2, v3, etc.
7. **WebSocket Support**: Real-time updates
8. **Background Jobs**: Async processing

## Contributing

Follow the development workflow and ensure:
- All tests pass
- Code is formatted
- Linting passes
- Documentation is updated
- Pre-commit hook succeeds

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Gin Framework](https://gin-gonic.com/)
- [pgx Documentation](https://github.com/jackc/pgx)
