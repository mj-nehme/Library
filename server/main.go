package main

import (
	"context"
	"library/api"
	"library/config"
	"library/db"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
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
	slog.Info("Configuration: ", cfg)

	// Initialize the database connection
	db := db.New()
	err = db.Connect(&cfg.Database)
	if err != nil {
		slog.Error("Error connecting to Database")
	}

	// Start the API server
	router := api.SetupRouter(db)
	server, err := api.StartServer(ctx, cfg.Server.Port, router)
	if err != nil {
		slog.Error("Failed to start the API server")
	}
	slog.Info("Started server %s...", server.Addr)
}
