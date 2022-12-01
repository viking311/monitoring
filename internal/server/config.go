package server

import (
	"os"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func ReadConfig() *ServerConfig {
	cfg := ServerConfig{}

	if err := env.Parse(&cfg); err != nil {
		os.Exit(1)
	}

	return &cfg
}
