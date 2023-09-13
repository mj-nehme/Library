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

	dbHost    = "172.21.0.5"
	dbPort    = 5432
	dbName    = "Library"
	dbUser    = "postgres"
	dbPass    = "postgres"
	dbSslMode = "disable"

	srvHost = "localhost"
	srvPort = 8090
)

func TestLoadConfigFromEnv(t *testing.T) {
	// Set the environment variables for the test
	os.Setenv("POSTGRES_HOST", dbHost)
	os.Setenv("POSTGRES_PORT", strconv.Itoa(dbPort))
	os.Setenv("POSTGRES_NAME", dbName)
	os.Setenv("POSTGRES_USERNAME", dbUser)
	os.Setenv("POSTGRES_PASSWORD", dbPass)
	os.Setenv("POSTGRES_SSL_MODE", dbSslMode)
	os.Setenv("SERVER_HOST", srvHost)
	os.Setenv("SERVER_PORT", strconv.Itoa(srvPort))

	// Clean up the environment variables after the test
	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_NAME")
		os.Unsetenv("POSTGRES_USERNAME")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_SSL_MODE")
		os.Unsetenv("SERVER_HOST")
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
	assert.Equal(t, srvHost, cfg.Server.Host)
	assert.Equal(t, srvPort, cfg.Server.Port)
}
