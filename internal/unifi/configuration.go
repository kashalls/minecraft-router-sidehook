package unifi

import (
	"github.com/caarlos0/env/v11"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

type UnifiConfig struct {
	Host               string `env:"UNIFI_HOST,notEmpty"`
	ApiKey             string `env:"UNIFI_API_KEY" envDefault:""`
	User               string `env:"UNIFI_USER" envDefault:""`
	Password           string `env:"UNIFI_PASS" envDefault:""`
	Site               string `env:"UNIFI_SITE" envDefault:"default"`
	ExternalController bool   `env:"UNIFI_EXTERNAL_CONTROLLER" envDefault:"false"`
	SkipTLSVerify      bool   `env:"UNIFI_SKIP_TLS_VERIFY" envDefault:"true"`

	IPv4ObjectName         string `env:"UNIFI_IPV4_OBJECT_NAME" envDefault:"Minecraft Router Block List v4"`
	IPv6ObjectName         string `env:"UNIFI_IPV6_OBJECT_NAME" envDefault:"Minecraft Router Block List v6"`
	VerifyObjects          bool   `env:"UNIFI_VERIFY_OBJECTS" envDefault:"true"`
	IPv4DefaultObjectValue string `env:"UNIFI_IPV4_DEFAULT_OBJECT_VALUE" envDefault:"255.255.255.255"`
	IPv6DefaultObjectValue string `env:"UNIFI_IPV6_DEFAULT_OBJECT_VALUE" envDefault:"ff02::1"`
}

var Config UnifiConfig

func InitConfig() UnifiConfig {
	Config := UnifiConfig{}
	if err := env.Parse(&Config); err != nil {
		log.Error("error reading configuration from environment", zap.Error(err))
	}
	return Config
}
