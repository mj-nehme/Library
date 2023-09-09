package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

var server *http.Server = nil

// StartServer starts the API server with the provided configuration.
func StartServer(ctx context.Context, port int, router *gin.Engine) error {
	addr := ":" + strconv.Itoa(port)
	slog.Info("Starting API server on %s...", addr)

	server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Failed to start the server: %v", err)
	}

	return nil
}

// ShutdownServer shuts down the API server gracefully.
func ShutdownServer(ctx context.Context) error {
	slog.Info("Shutting down API server gracefully...")

	// Attempt to shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		// If the Shutdown method returns an error, forcefully close the server
		return fmt.Errorf("Failed to start the server: %v", err)
		//server.Close()
	}

	return nil
}
