package discord

import (
	"github.com/caarlos0/env/v11"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

type DiscordConfig struct {
	WebhookURL           string `env:"WEBHOOK_URL" envDefault:""`
	WebhookURLIsTemplate bool   `env:"WEBHOOK_URL_IS_TEMPLATE" envDefault:"false"`
	DiscordTemplate      string `env:"DISCORD_TEMPLATE" envDefault:""`
}

func InitConfig() DiscordConfig {
	cfg := DiscordConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Error("error reading configuration from environment", zap.Error(err))
	}
	return cfg
}
