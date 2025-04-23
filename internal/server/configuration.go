package server

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

type Config struct {
	ServerHost         string        `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort         int           `env:"SERVER_PORT" envDefault:"8888"`
	ServerReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`

	HealthHost         string        `env:"HEALTH_HOST" envDefault:"localhost"`
	HealthPort         int           `env:"HEALTH_PORT" envDefault:"8080"`
	HealthReadTimeout  time.Duration `env:"HEALTH_READ_TIMEOUT"`
	HealthWriteTimeout time.Duration `env:"HEALTH_WRITE_TIMEOUT"`
}

func InitConfig() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Error("error reading configuration from environment", zap.Error(err))
	}
	return cfg
}
