package config

// Config holds the configuration settings for the Library application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig holds the database configuration settings
type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USERNAME"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_NAME"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

// ServerConfig holds the server configuration settings
type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port int    `env:"SERVER_PORT"`
	//RootDir string `env:"ROOT_DIRECTORY"`
}
