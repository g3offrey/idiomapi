package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/g3offrey/idiomapi/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database wraps the pgx connection pool
type Database struct {
	Pool   *pgxpool.Pool
	logger *slog.Logger
}

// New creates a new Database instance with a connection pool
func New(ctx context.Context, cfg *config.DatabaseConfig, logger *slog.Logger) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	if cfg.MaxOpenConns > 0 && cfg.MaxOpenConns <= 2147483647 {
		poolConfig.MaxConns = int32(cfg.MaxOpenConns) // #nosec G115
	}
	if cfg.MaxIdleConns > 0 && cfg.MaxIdleConns <= 2147483647 {
		poolConfig.MinConns = int32(cfg.MaxIdleConns) // #nosec G115
	}
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName)

	return &Database{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	db.logger.Info("closing database connection")
	db.Pool.Close()
}

// Health checks the database connection health
func (db *Database) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
