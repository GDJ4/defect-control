package config

import (
	"time"

	"github.com/caarlos0/env/v10"
)

// Config aggregates all runtime configuration knobs that can be overridden via environment variables.
type Config struct {
	AppName string `env:"APP_NAME" envDefault:"defect-tracker"`
	AppEnv  string `env:"APP_ENV" envDefault:"development"`

	Server struct {
		Host         string        `env:"SERVER_HOST" envDefault:"0.0.0.0"`
		Port         int           `env:"SERVER_PORT" envDefault:"8080"`
		ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
		IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	}

	Database struct {
		DSN string `env:"DATABASE_DSN" envDefault:"postgres://user:password@localhost:5432/defects?sslmode=disable"`
	}

	Storage struct {
		Driver string `env:"STORAGE_DRIVER" envDefault:"local"`
		Path   string `env:"STORAGE_PATH" envDefault:"storage/uploads"`
		S3     struct {
			Endpoint   string        `env:"STORAGE_S3_ENDPOINT"`
			Bucket     string        `env:"STORAGE_S3_BUCKET"`
			AccessKey  string        `env:"STORAGE_S3_ACCESS_KEY"`
			SecretKey  string        `env:"STORAGE_S3_SECRET_KEY"`
			Region     string        `env:"STORAGE_S3_REGION" envDefault:"us-east-1"`
			UseSSL     bool          `env:"STORAGE_S3_USE_SSL" envDefault:"false"`
			PresignTTL time.Duration `env:"STORAGE_S3_PRESIGN_TTL" envDefault:"15m"`
		}
	}

	Auth struct {
		Secret     string        `env:"JWT_SECRET,required"`
		AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" envDefault:"1h"`
		RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"` // default 30 days
	}
}

// Load parses environment variables into Config.
func Load() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
