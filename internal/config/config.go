package config

import (
	"github.com/caarlos0/env/v10"
	"time"
)

// Config dsn format: "user=user password=pass host=host port=port dbname=name"
type Config struct {
	APIKey         string        `env:"APOD_API,required"`
	DataBaseDSN    string        `env:"DATABASE_DSN,required"`
	MinioEndpoint  string        `env:"M_ENDPOINT,required"`
	MinioUser      string        `env:"M_USER,required"`
	MinioPass      string        `env:"M_PASS,required"`
	MinioBucket    string        `env:"M_BUCKET,required"`
	MinioPort      string        `env:"M_PORT,required"`
	ServerAddress  string        `env:"SERVER_ADDRESS" envDefault:"localhost"`
	ServerPort     string        `env:"SERVER_PORT" envDefault:"8080"`
	WorkerInterval time.Duration `env:"WORKER_INTERVAL" envDefault:"24h"`
	WorkerRetries  int           `env:"WORKER_RETRIES" envDefault:"5"`
}

func MustLoad() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil || cfg.APIKey == "" {
		panic(err)
	}
	return &cfg
}
