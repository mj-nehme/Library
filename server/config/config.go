package config

import (
	"context"

	"log"

	"github.com/caarlos0/env"
)

func Load(ctx context.Context) (Config, error) {
	var dbConfig DatabaseConfig
	var serverConfig ServerConfig

	if err := env.Parse(&dbConfig); err != nil {
		log.Fatal("Error parsing database config:", err)
	}

	if err := env.Parse(&serverConfig); err != nil {
		log.Fatal("Error parsing server config:", err)
	}

	return Config{Database: dbConfig, Server: serverConfig}, nil
}
