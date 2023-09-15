package config

import (
	"context"
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

	srvHost = "localhost"
	srvPort = 8090
)

func TestLoadConfigFromEnv(t *testing.T) {
	err := addEnvirnomentVariables()
	assert.NoError(t, err)

	// Call the function to load the configuration from environment variables
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout*time.Second))
	defer cancel()
	cfg, err := parseEnvironmentVariables(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check individual fields in the Config struct
	assert.Equal(t, dbHost, cfg.Database.Host)
	assert.Equal(t, dbPort, cfg.Database.Port)
	assert.Equal(t, dbName, cfg.Database.Name)
	assert.Equal(t, dbUser, cfg.Database.Username)
	assert.Equal(t, dbPass, cfg.Database.Password)
	assert.Equal(t, dbSslMode, cfg.Database.SSLMode)
	assert.Equal(t, srvHost, cfg.Server.Host)
	assert.Equal(t, srvPort, cfg.Server.Port)
}
