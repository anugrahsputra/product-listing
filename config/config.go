package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/op/go-logging"
)

var configLog = logging.MustGetLogger("config")

type Config struct {
	Port       string `env:"PORT" env-default:"8080"`
	DBHost     string `env:"DB_HOST" env-required:"true"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-required:"true"`
	DBPassword string `env:"DB_PASSWORD" env-required:"true"`
	DBName     string `env:"DB_NAME" env-required:"true"`
}

func Load() *Config {
	var cfg Config

	// Try to load from .env file, but don't fail if it doesn't exist
	_ = cleanenv.ReadConfig(".env", &cfg)

	// Read from environment variables (overrides .env)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		configLog.Warningf("Warning: Failed to read environment variables: %v", err)
	}

	return &cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
