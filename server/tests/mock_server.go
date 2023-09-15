package tests

import (
	"context"
	"fmt"
	"library/api"
	"library/config"
	"library/db"
	"library/models"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	host            = "http://localhost"
	apiPath         = "/api/v1"
	apiVersion      = "1.0.0"
	healthCheckPath = "/health"
	testPort        = 8088
	contextTimeout  = 10
	testDatabase    = "Library"
)

var ServerAddress = ""

func init() {
	gin.SetMode(gin.TestMode)
}

func UpdateHostAddress(host string, port int) {
	ServerAddress = host + ":" + strconv.Itoa(port)
}

func SetupMockServer() (context.Context, *gin.Engine, db.Database) {
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
	cfg.Database.Name = testDatabase
	err = db.Connect(&cfg.Database)
	if err != nil {
		slog.Error("Error connecting to Database")
	}
	time.Sleep(time.Second)

	// Start the API server
	router := api.SetupRouter(db)

	// Choose some arbitrary port for that consecutive tests
	// might lead to ports in a CLOSE_WAIT status
	//port := GenerateArbitraryPort()
	port := cfg.Server.Port

	// Start the server
	errChan := make(chan error)
	// Start the server in a goroutine
	go func() {
		err := api.StartServer(ctx, cfg.Server.Port, router)
		errChan <- err // Send the error back through the channel
	}()

	// Check for errors
	select {
	case err := <-errChan:
		if err != nil {
			slog.Error("Failed to start the server: %v", err)
		}
	case <-time.After(time.Second * 5): // Timeout after 5 seconds
		slog.Info("Server started successfully")
	}

	UpdateHostAddress(host, port)
	cfg.Server.Port = port

	// Wait for the server to be ready
	if err := waitForServerReady(ServerAddress); err != nil {
		slog.Error(err.Error())
	}

	return ctx, router, db
}

func TearDownMockServer(ctx context.Context, db db.Database) {
	// Shut down the server gracefully
	err := api.ShutdownServer(ctx)
	if err != nil {
		slog.Error(err.Error())
	}

	err = db.DB.Migrator().DropTable(&models.Book{})
	if err != nil {
		slog.Error(err.Error())
	}
}

func waitForServerReady(address string) error {
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		// Send a request to the readiness endpoint
		resp, err := http.Get(address + healthCheckPath)
		if err == nil && resp.StatusCode == http.StatusOK {
			// Server is ready, return
			return nil
		}

		// Wait for a short duration before the next attempt
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("server is not ready")
}

func GenerateArbitraryPort() int {
	portRange := 65535 - 1024
	port := rand.Intn(portRange)
	port += 1024

	return port
}
