package agent

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultAddress        = "localhost:8080"
	DefaultReportInterval = 10 * time.Second
	DefaultPollInterval   = 2 * time.Second
)

type AgentConfig struct {
	Address        *string        `env:"ADDRESS"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
}

var Config AgentConfig

func init() {
	addressFlag := flag.String("a", DefaultAddress, "address to send metrics")
	reportInterval := flag.Duration("r", DefaultReportInterval, "how often send report to server")
	pollInterval := flag.Duration("p", DefaultPollInterval, "how often update metrics")
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
