package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	content := `
[server]
host = "localhost"
port = 8080
read_timeout = "15s"
write_timeout = "15s"
idle_timeout = "60s"

[database]
host = "localhost"
port = 5432
user = "testuser"
password = "testpass"
dbname = "testdb"
sslmode = "disable"
max_open_conns = 25
max_idle_conns = 25
conn_max_lifetime = "5m"

[logging]
level = "info"
format = "json"
add_source = false
`
	tmpfile, err := os.CreateTemp("", "config-*.toml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	assert.NoError(t, err)
	tmpfile.Close()

	// Test Load
	cfg, err := Load(tmpfile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify server config
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 15*time.Second, cfg.Server.ReadTimeout)

	// Verify database config
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testdb", cfg.Database.DBName)

	// Verify logging config
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestServerConfig_Address(t *testing.T) {
	cfg := ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}
	assert.Equal(t, "0.0.0.0:8080", cfg.Address())
}

func TestDatabaseConfig_DSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}

func TestLoad_InvalidFile(t *testing.T) {
	_, err := Load("nonexistent.toml")
	assert.Error(t, err)
}
