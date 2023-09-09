package config

import (
	"context"
	"library/tools"

	"github.com/spf13/viper"
)

var (
	configFileName      = "config"
	configFileExtension = "toml"
)

func Load(ctx context.Context) (Config, error) {
	configFilePath, err := tools.SearchRootDirectory()
	if err != nil {
		return Config{}, err
	}
	viper.AddConfigPath(configFilePath)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	// Unmarshal DatabaseConfig
	if err = viper.UnmarshalKey("Database", &config.Database); err != nil {
		return Config{}, err
	}

	// Unmarshal ServerConfig
	if err = viper.UnmarshalKey("Server", &config.Server); err != nil {
		return Config{}, err
	}
	return config, nil
}
