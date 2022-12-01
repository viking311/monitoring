package agent

import (
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func ReadConfig() *AgentConfig {
	cfg := AgentConfig{}

	if err := env.Parse(&cfg); err != nil {
		os.Exit(1)
	}

	return &cfg
}
