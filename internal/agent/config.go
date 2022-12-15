package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/viking311/monitoring/internal/logger"
)

const (
	DefaultAddress        = "localhost:8080"
	DefaultReportInterval = 10 * time.Second
	DefaultPollInterval   = 2 * time.Second
	DefaultHashKey        = ""
)

type AgentConfig struct {
	Address        *string        `env:"ADDRESS"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
	HashKey        *string        `env:"KEY"`
}

var Config AgentConfig

func init() {
	logger.Logger.Info("start reading configuration")

	logger.Logger.Debug("reading flags")
	addressFlag := flag.String("a", DefaultAddress, "address to send metrics")
	reportInterval := flag.Duration("r", DefaultReportInterval, "how often send report to server")
	pollInterval := flag.Duration("p", DefaultPollInterval, "how often update metrics")
	hashKey := flag.String("k", DefaultHashKey, "hash key")
	flag.Parse()

	logger.Logger.Debug("reading enviroments")
	if err := env.Parse(&Config); err != nil {
		logger.Logger.Fatal(err)
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

	if Config.HashKey == nil {
		Config.HashKey = hashKey
	}

	logger.Logger.Info("finish reading configuration")
}
