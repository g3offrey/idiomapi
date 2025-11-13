package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/g3offrey/idiomapi/internal/config"
	"github.com/g3offrey/idiomapi/internal/database"
	"github.com/g3offrey/idiomapi/internal/handler"
	"github.com/g3offrey/idiomapi/internal/middleware"
	"github.com/g3offrey/idiomapi/internal/repository"
	"github.com/g3offrey/idiomapi/internal/service"
	"github.com/g3offrey/idiomapi/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/config.toml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging)
	log.Info("starting application",
		"config", *configPath,
		"server_address", cfg.Server.Address())

	// Initialize database
	ctx := context.Background()
	db, err := database.New(ctx, &cfg.Database, log)
	if err != nil {
		log.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(db.Pool)

	// Initialize services
	todoService := service.NewTodoService(todoRepo, log)

	// Initialize handlers
	todoHandler := handler.NewTodoHandler(todoService)
	healthHandler := handler.NewHealthHandler(db)

	// Setup Gin
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logger(log))

	// Setup routes
	setupRoutes(router, todoHandler, healthHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info("server starting", "address", cfg.Server.Address())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server stopped")
}

// setupRoutes configures all API routes
func setupRoutes(router *gin.Engine, todoHandler *handler.TodoHandler, healthHandler *handler.HealthHandler) {
	// Health check
	router.GET("/health", healthHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	todos := v1.Group("/todos")
	todos.POST("", todoHandler.CreateTodo)
	todos.GET("", todoHandler.ListTodos)
	todos.GET("/:id", todoHandler.GetTodo)
	todos.PUT("/:id", todoHandler.UpdateTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)
}
