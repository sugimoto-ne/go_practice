package config

import (
	"github.com/caarlos0/env/v7"
)

type Config struct {
	Env                     string `env:"APP_ENV" envDefault:"dev"`
	Port                    int    `env:"PORT" envDefault:"8080"`
	ReadHeaderTimeoutSecond int    `env:"READ_HEADER_TIMEOUT_SECOND" envDefault:"20"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
