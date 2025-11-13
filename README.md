# idiomapi

A production-ready REST API for a TODO application built with Go, following Go best practices and idioms.

## Features

- 🚀 **RESTful API** - Clean REST API design with proper HTTP methods and status codes
- 🔒 **Type-safe Database** - PostgreSQL with pgx/v5 for high-performance database operations
- 📝 **Structured Logging** - Configurable slog-based logging with JSON/text formats
- ⚙️ **Configuration Management** - TOML-based configuration using cleanenv
- 🧪 **Comprehensive Testing** - Unit tests using testify/assert
- 🔍 **Code Quality** - GolangCI-lint integration with pre-commit hooks
- 🐳 **Docker Support** - Docker Compose for easy local development
- 📊 **Health Checks** - Built-in health check endpoints
- 🎯 **Clean Architecture** - Well-organized project structure with separation of concerns

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [pgx](https://github.com/jackc/pgx) - Pure Go PostgreSQL driver
- **Configuration**: [cleanenv](https://github.com/ilyakaznacheev/cleanenv) - TOML/YAML/ENV configuration
- **Logging**: [slog](https://pkg.go.dev/log/slog) - Structured logging (Go 1.21+)
- **Testing**: [testify](https://github.com/stretchr/testify) - Testing toolkit with assertions
- **Linting**: [GolangCI-lint](https://golangci-lint.run/) - Fast Go linters aggregator

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entrypoint
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection and setup
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # Data models and DTOs
│   ├── repository/       # Data access layer
│   └── service/          # Business logic layer
├── pkg/
│   └── logger/           # Logging utilities
├── migrations/           # Database migrations
├── configs/              # Configuration files
├── scripts/              # Build and deployment scripts
├── Makefile             # Build automation
└── docker-compose.yml   # Docker services
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker & Docker Compose (optional, for local development)
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/g3offrey/idiomapi.git
cd idiomapi
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Start PostgreSQL (using Docker)

```bash
make docker-up
```

Or manually:
```bash
docker-compose up -d
```

### 4. Create database (if not using docker-compose)

```bash
make db-create
```

### 5. Run migrations

```bash
make migrate-up
```

### 6. Start the application

```bash
make run
```

The API will be available at `http://localhost:8080`

## Configuration

Configuration is managed through TOML files in the `configs/` directory. The default configuration file is `configs/config.toml`.

```toml
[server]
host = "0.0.0.0"
port = 8080
read_timeout = "15s"
write_timeout = "15s"
idle_timeout = "60s"

[database]
host = "localhost"
port = 5432
user = "postgres"
password = "postgres"
dbname = "tododb"
sslmode = "disable"
max_open_conns = 25
max_idle_conns = 25
conn_max_lifetime = "5m"

[logging]
level = "info"  # debug, info, warn, error
format = "json" # json, text
add_source = false
```

You can override the config file path using the `-config` flag:

```bash
go run cmd/api/main.go -config /path/to/config.toml
```

## API Endpoints

### Health Check

```
GET /health
```

### Todos

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/todos` | Create a new todo |
| GET | `/api/v1/todos` | List all todos (with pagination) |
| GET | `/api/v1/todos/:id` | Get a specific todo |
| PUT | `/api/v1/todos/:id` | Update a todo |
| DELETE | `/api/v1/todos/:id` | Delete a todo |

### Example Requests

**Create a todo:**
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": false
  }'
```

**List todos:**
```bash
curl http://localhost:8080/api/v1/todos?page=1&page_size=10
```

**Get a todo:**
```bash
curl http://localhost:8080/api/v1/todos/1
```

**Update a todo:**
```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'
```

**Delete a todo:**
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

**Filter by completion status:**
```bash
curl http://localhost:8080/api/v1/todos?completed=true
```

## Development

### Build

```bash
make build
```

### Run

```bash
make run
```

### Test

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

### Linting

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

### Run all CI checks locally

```bash
make ci
```

## Git Hooks

Install pre-commit hooks to run CI checks locally before committing:

```bash
make install-hooks
```

This will:
- Check code formatting
- Run go vet
- Run golangci-lint
- Run tests

## Database Migrations

Create the database:
```bash
make db-create
```

Run migrations:
```bash
make migrate-up
```

Rollback migrations:
```bash
make migrate-down
```

Drop the database:
```bash
make db-drop
```

## Docker

Start all services:
```bash
make docker-up
```

Stop all services:
```bash
make docker-down
```

## Environment Variables

You can override configuration values using environment variables (when using cleanenv):

```bash
export SERVER_PORT=9090
export DATABASE_HOST=remote-db.example.com
export LOGGING_LEVEL=debug
```

## Code Quality

This project follows Go best practices and idioms:

- ✅ Clean architecture with clear separation of concerns
- ✅ Dependency injection for testability
- ✅ Context propagation for cancellation and timeouts
- ✅ Structured logging with levels
- ✅ Proper error handling and wrapping
- ✅ Connection pooling for database
- ✅ Graceful shutdown
- ✅ Input validation
- ✅ Comprehensive testing
- ✅ Pre-commit hooks for code quality

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Maintainers

- [@g3offrey](https://github.com/g3offrey)