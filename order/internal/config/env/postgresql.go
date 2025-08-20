package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresqlEnvConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	Database string `env:"POSTGRES_DB,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE,required"`
}

type postgresqlConfig struct {
	raw postgresqlEnvConfig
}

func NewPostgreSQLConfig() (*postgresqlConfig, error) {
	var raw postgresqlEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &postgresqlConfig{raw: raw}, nil
}

func (cfg *postgresqlConfig) URI() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.SSLMode,
	)
}

func (cfg *postgresqlConfig) DatabaseName() string {
	return cfg.raw.Database
}
