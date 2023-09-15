package db

import (
	"context"
	"library/config"
	"testing"
	"time"

	"log/slog"

	"github.com/stretchr/testify/assert"
)

const (
	contextTimeout = 10
	tempDatabase   = "Library"
)

func SetupTest(t *testing.T) Database {
	slog.Info("Setting up Test..")

	// Load config
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout*time.Second))
	defer cancel()
	cfg, err := config.Load(ctx)
	assert.NoError(t, err)
	assert.NotZero(t, cfg.Database.Port)
	assert.NotZero(t, cfg.Server.Port)

	// Connect to DB
	db := New()
	cfg.Database.Name = tempDatabase
	err = db.Connect(&cfg.Database)
	assert.NoError(t, err)
	time.Sleep(time.Second)

	slog.Info("Setup'ed test successfully..")

	return db
}
