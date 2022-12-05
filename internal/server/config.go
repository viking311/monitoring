package server

import (
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func ReadConfig() *ServerConfig {
	cfg := ServerConfig{}

	if err := env.Parse(&cfg); err != nil {
		os.Exit(1)
	}

	return &cfg
}
