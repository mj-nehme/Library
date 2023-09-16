//	@title		Library API
//	@version	2.0

package main

import (
	"context"
	"library/api"
	"library/config"
	"library/db"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	contextTimeout = 10
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	// Load config
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout*time.Second))
	defer cancel()
	cfg, err := config.Load(ctx)
	if err != nil {
		slog.Error("Error loading config file")
	}
	slog.Info("loaded configuration successfully.", "Configuration", cfg)

	// Initialize the database connection
	db := db.New()
	err = db.Connect(&cfg.Database)
	if err != nil {
		slog.Error("Error connecting to Database")
	}

	time.Sleep(time.Second)

	// Start the API server
	router := api.SetupRouter(db)
	err = api.StartServer(ctx, cfg.Server.Port, router)
	if err != nil {
		slog.Error("Failed to start the API server")
	}
}
