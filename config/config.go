package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Config holds the configuration for the application.
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig holds the configuration for the database.
type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	Name     string `env:"DB_NAME,required"`
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
