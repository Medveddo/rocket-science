package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT,required"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type httpConfig struct {
	raw httpEnvConfig
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpConfig{raw: raw}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *httpConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *httpConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *httpConfig) ReadHeaderTimeout() time.Duration {
	return cfg.raw.ReadHeaderTimeout
}

func (cfg *httpConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}