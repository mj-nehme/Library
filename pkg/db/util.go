package db

import (
	"context"
	"library/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const (
	contextTimeout = 10
)

func SetupTest(t *testing.T) Database {
	slog.Info("Setting up Test..")

	// Load config
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout*time.Second))
	defer cancel()
	cfg, err := config.Load(ctx)
	assert.NoError(t, err)

	// Connect to DB
	db := Database{DB: &gorm.DB{}}
	err = db.Connect(&cfg.Database)
	assert.NoError(t, err)

	slog.Info("Setup'ed test successfully..")

	return db
}
