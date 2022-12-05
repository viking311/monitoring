package agent

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	DEFAULT_ADDRESS         = "localhost:8080"
	DEFAULT_REPORT_INTERVAL = 10 * time.Second
	DEFAULT_POLL_INTERVAL   = 2 * time.Second
)

type AgentConfig struct {
	Address        *string        `env:"ADDRESS"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
}

var Config AgentConfig

func init() {
	addressFlag := flag.String("a", DEFAULT_ADDRESS, "address to send metrics")
	reportInterval := flag.Duration("r", DEFAULT_REPORT_INTERVAL, "how often send report to server")
	pollInterval := flag.Duration("p", DEFAULT_POLL_INTERVAL, "how often update metrics")
	flag.Parse()

	if err := env.Parse(&Config); err != nil {
		log.Fatal(err)
	}

	if Config.Address == nil {
		Config.Address = addressFlag
	}

	if Config.ReportInterval == nil {
		Config.ReportInterval = reportInterval
	}

	if Config.PollInterval == nil {
		Config.PollInterval = pollInterval
	}
}
