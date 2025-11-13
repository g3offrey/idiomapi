package logger

import (
	"log/slog"
	"testing"

	"github.com/g3offrey/idiomapi/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.LoggingConfig
		level  slog.Level
		format string
	}{
		{
			name: "json debug level",
			cfg: config.LoggingConfig{
				Level:     "debug",
				Format:    "json",
				AddSource: false,
			},
			level:  slog.LevelDebug,
			format: "json",
		},
		{
			name: "text info level",
			cfg: config.LoggingConfig{
				Level:     "info",
				Format:    "text",
				AddSource: true,
			},
			level:  slog.LevelInfo,
			format: "text",
		},
		{
			name: "warn level",
			cfg: config.LoggingConfig{
				Level:  "warn",
				Format: "json",
			},
			level: slog.LevelWarn,
		},
		{
			name: "error level",
			cfg: config.LoggingConfig{
				Level:  "error",
				Format: "json",
			},
			level: slog.LevelError,
		},
		{
			name: "default level on invalid",
			cfg: config.LoggingConfig{
				Level:  "invalid",
				Format: "json",
			},
			level: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.cfg)
			assert.NotNil(t, logger)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"invalid", slog.LevelInfo},
		{"", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
