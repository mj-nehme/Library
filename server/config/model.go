package config

// Config holds the configuration settings for the Library application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig holds the database configuration settings
type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
}

// ServerConfig holds the server configuration settings
type ServerConfig struct {
	Port int `mapstructure:"SERVER_PORT"`
	//RootDir string `mapstructure:"ROOT_DIRECTORY"`
}
