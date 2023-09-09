package core

import (
	"fmt"
	"library/config"
	"log"

	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDB() Database {
	return Database{DB: &gorm.DB{}}
}

// InitDB initializes the database connection pool
func (db *Database) Connect(cfg *config.DatabaseConfig) error {
	connectionString := buildDatabaseConnectionString(cfg)
	slog.Debug("Connection String:", connectionString)

	var err error
	db.DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database. %v", err)
	}

	return nil
}

func buildDatabaseConnectionString(cfg *config.DatabaseConfig) string {
	connectionString := fmt.Sprintf("host=%v port=%d dbname=%v sslmode=%v", cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	if cfg.User != "" {
		// If the user is not provided in the config, assume there is no user authentication
		connectionString += fmt.Sprintf(" user=%v password=%v", cfg.User, cfg.Password)
	}

	return connectionString
}
