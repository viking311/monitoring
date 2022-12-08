package main

import (
	"github.com/viking311/monitoring/internal/agent"
	"github.com/viking311/monitoring/internal/entity"
)

func main() {
	collector := entity.NewCollector("http://"+*agent.Config.Address, *agent.Config.PollInterval, *agent.Config.ReportInterval)

	collector.Do()

}
