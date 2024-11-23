package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Config holds the configuration for the application.
type Config struct {
	Env          string `env:"ENV,default=development"`
	Database     DatabaseConfig
	Server       ServerConfig
	LocalStorage LocalStorageConfig
}

// DatabaseConfig holds the configuration for the database.
type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	Name     string `env:"DB_NAME,required"`
}

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	Port int `env:"SERVER_PORT,default=8080"`
}

type LocalStorageConfig struct {
	UseLocal bool   `env:"USE_LOCAL_STORAGE"`
	Dir      string `env:"LOCAL_STORAGE_DIR"`
}

// NewConfig creates an instance of Config.
func NewConfig(env string) (*Config, error) {
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "error decoding env")
	}

	return &config, nil
}
