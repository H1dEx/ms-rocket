package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	Port string `env:"HTTP_PORT,required"`
	Host string `env:"HTTP_HOST,required"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &orderHTTPConfig{raw: raw}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
