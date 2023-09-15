package config

import (
	"context"
	"fmt"
	"library/tools"
	"os"
	"strconv"

	"log"

	"github.com/spf13/viper"

	"github.com/caarlos0/env"
)

var (
	configFileName      = "config"
	configFileExtension = "env"
)

func Load(ctx context.Context) (Config, error) {
	err := addEnvirnomentVariables()
	if err != nil {
		return Config{}, err
	}
	return parseEnvironmentVariables(ctx)
}

func parseEnvironmentVariables(ctx context.Context) (Config, error) {
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

func addEnvirnomentVariables() error {
	config, err := loadConfigFile()
	if err != nil {
		return fmt.Errorf("cannot load config file: %s", err)
	}
	err = addEnvirnomentVariablesFromFile(config)
	if err != nil {
		return fmt.Errorf("cannot add environment variables: %s", err)
	}

	return nil
}

func loadConfigFile() (config Config, err error) {
	configFilePath, err := tools.SearchRootDirectory()
	if err != nil {
		log.Fatal("cannot find root directory:", err)
	}
	viper.AddConfigPath(configFilePath + "/../")
	viper.AddConfigPath(configFilePath)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	fmt.Println("viper AllKeys:", viper.AllKeys())
	fmt.Println("viper AllSettings:", viper.AllSettings())

	var databaseConfig DatabaseConfig
	var serverConfig ServerConfig

	var POSTGRES_HOST string
	err = viper.UnmarshalKey("POSTGRES_HOST", &POSTGRES_HOST)
	if err != nil {
		return Config{}, err
	}
	var POSTGRES_PORT int
	err = viper.UnmarshalKey("POSTGRES_PORT", &POSTGRES_PORT)
	if err != nil {
		return Config{}, err
	}
	var POSTGRES_USERNAME string
	err = viper.UnmarshalKey("POSTGRES_USERNAME", &POSTGRES_USERNAME)
	if err != nil {
		return Config{}, err
	}
	var POSTGRES_PASSWORD string
	err = viper.UnmarshalKey("POSTGRES_PASSWORD", &POSTGRES_PASSWORD)
	if err != nil {
		return Config{}, err
	}
	var POSTGRES_NAME string
	err = viper.UnmarshalKey("POSTGRES_NAME", &POSTGRES_NAME)
	if err != nil {
		return Config{}, err
	}
	var POSTGRES_SSL_MODE string
	err = viper.UnmarshalKey("POSTGRES_SSL_MODE", &POSTGRES_SSL_MODE)
	if err != nil {
		return Config{}, err
	}
	var SERVER_HOST string
	err = viper.UnmarshalKey("SERVER_HOST", &SERVER_HOST)
	if err != nil {
		return Config{}, err
	}
	var SERVER_PORT int
	err = viper.UnmarshalKey("SERVER_PORT", &SERVER_PORT)
	if err != nil {
		return Config{}, err
	}

	databaseConfig = DatabaseConfig{POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USERNAME, POSTGRES_PASSWORD, POSTGRES_NAME, POSTGRES_SSL_MODE}
	serverConfig = ServerConfig{SERVER_HOST, SERVER_PORT}

	return Config{databaseConfig, serverConfig}, nil
}

func addEnvirnomentVariablesFromFile(config Config) error {
	// Set the environment variables for the test
	err := os.Setenv("POSTGRES_HOST", config.Database.Host)
	if err != nil {
		return err
	}
	err = os.Setenv("POSTGRES_PORT", strconv.Itoa(config.Database.Port))
	if err != nil {
		return err
	}
	err = os.Setenv("POSTGRES_NAME", config.Database.Name)
	if err != nil {
		return err
	}
	err = os.Setenv("POSTGRES_USERNAME", config.Database.Username)
	if err != nil {
		return err
	}
	err = os.Setenv("POSTGRES_PASSWORD", config.Database.Password)
	if err != nil {
		return err
	}
	err = os.Setenv("POSTGRES_SSL_MODE", config.Database.SSLMode)
	if err != nil {
		return err
	}
	err = os.Setenv("SERVER_HOST", config.Server.Host)
	if err != nil {
		return err
	}
	err = os.Setenv("SERVER_PORT", strconv.Itoa(config.Server.Port))
	if err != nil {
		return err
	}

	return nil
}
