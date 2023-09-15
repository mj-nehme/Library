package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// StartServer starts the API server with the provided configuration.
func StartServer(ctx context.Context, port int, router *gin.Engine) (*http.Server, error) {
	portSrt := strconv.Itoa(port)
	addr := ":" + portSrt
	slog.Info("Starting API server on port.", "port", portSrt)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return nil, fmt.Errorf("failed to start the server: %v", err)
	}

	return server, nil
}

// ShutdownServer shuts down the API server gracefully.

func ShutdownServer(ctx context.Context, server *http.Server) error {
	slog.Info("Shutting down API server gracefully...")

	// Attempt to shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		// If the Shutdown method returns an error, forcefully close the server
		return fmt.Errorf("failed to start the server: %v", err)
		//server.Close()
	}

	return nil
}
