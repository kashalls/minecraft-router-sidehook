package configuration

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type DiscordConfig struct {
	Webhook         string `env:"WEBHOOK" envDefault:""`
	WebhookTemplate string `env:"WEBHOOK_TEMPLATE" envDefault:""`
	Token           string `env:"DISCORD_TOKEN" envDefault:""`
}

var Config DiscordConfig

func Init() DiscordConfig {
	Config = DiscordConfig{}
	if err := env.Parse(&Config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return Config
}
