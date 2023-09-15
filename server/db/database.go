package db

import (
	"errors"
	"fmt"
	"library/config"
	"library/models"
	"log"

	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New() Database {
	return Database{DB: &gorm.DB{}}
}

// InitDB initializes the database connection pool
func (db *Database) Connect(cfg *config.DatabaseConfig) error {
	slog.Info("Connecting to database.", "Host", cfg.Host, "Port", cfg.Port)
	connectionString := buildDatabaseConnectionString(cfg)
	slog.Debug("Built database connection string.", "Connection String", connectionString)

	var err error
	db.DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database. %v", err)
	}

	err = db.DB.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connected to database successfully..")
	return nil
}

func buildDatabaseConnectionString(cfg *config.DatabaseConfig) string {
	connectionString := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	if cfg.User != "" {
		// If the user is not provided in the config, assume there is no user authentication
		connectionString += fmt.Sprintf(" user=%s password=%s", cfg.User, cfg.Password)
	}

	return connectionString
}

// Teardown cleans up the database after testing
func (db *Database) Teardown() error {
	if db.DB == nil {
		return errors.New("database is pointing to nil")
	}

	return db.DB.Migrator().DropTable(&models.Book{})
}
