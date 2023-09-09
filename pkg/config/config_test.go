package config

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	contextTimeout = 10

	dbHost    = "localhost"
	dbPort    = 5432
	dbName    = "Library"
	dbUser    = "postgres"
	dbPass    = "postgres"
	dbSslMode = "disable"
	svPort    = 8090
)

func TestLoadConfigFromFile(t *testing.T) {
	// Call the function to load the configuration from environment variables
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout*time.Second))
	defer cancel()
	cfg, err := Load(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check the loaded values
	assert.Equal(t, dbHost, cfg.Database.Host)
	assert.Equal(t, dbPort, cfg.Database.Port)
	assert.Equal(t, dbName, cfg.Database.Name)
	assert.Equal(t, dbUser, cfg.Database.User)
	assert.Equal(t, dbPass, cfg.Database.Password)
	assert.Equal(t, dbSslMode, cfg.Database.SSLMode)
	assert.Equal(t, svPort, cfg.Server.Port)
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set the environment variables for the test
	os.Setenv("DB_HOST", dbHost)
	os.Setenv("DB_PORT", strconv.Itoa(dbPort))
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_USERNAME", dbUser)
	os.Setenv("DB_PASSWORD", dbPass)
	os.Setenv("DB_SSL_MODE", dbSslMode)
	os.Setenv("SERVER_PORT", strconv.Itoa(svPort))

	// Clean up the environment variables after the test
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSL_MODE")
		os.Unsetenv("SERVER_PORT")
	}()

	// Call the function to load the configuration from environment variables
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout*time.Second))
	defer cancel()
	cfg, err := Load(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check individual fields in the Config struct
	assert.Equal(t, dbHost, cfg.Database.Host)
	assert.Equal(t, dbPort, cfg.Database.Port)
	assert.Equal(t, dbName, cfg.Database.Name)
	assert.Equal(t, dbUser, cfg.Database.User)
	assert.Equal(t, dbPass, cfg.Database.Password)
	assert.Equal(t, dbSslMode, cfg.Database.SSLMode)
	assert.Equal(t, svPort, cfg.Server.Port)
}
