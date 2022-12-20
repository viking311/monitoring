package main

import (
	"github.com/viking311/monitoring/internal/agent"
	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/logger"
)

func main() {
	logger.Debug("start agent")
	defer logger.Debug("stop agent")

	collector := entity.NewCollector("http://"+*agent.Config.Address, *agent.Config.PollInterval, *agent.Config.ReportInterval, *agent.Config.HashKey)

	collector.Do()
}
