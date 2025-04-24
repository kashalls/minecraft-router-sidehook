package discord

import (
	"github.com/caarlos0/env/v11"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

type DiscordConfig struct {
	Webhook         string `env:"WEBHOOK" envDefault:""`
	WebhookTemplate string `env:"WEBHOOK_TEMPLATE" envDefault:""`
}

func InitConfig() DiscordConfig {
	cfg := DiscordConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Error("error reading configuration from environment", zap.Error(err))
	}
	return cfg
}
