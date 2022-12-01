package main

import (
	"github.com/viking311/monitoring/internal/agent"
	"github.com/viking311/monitoring/internal/entity"
)

func main() {
	cfg := agent.ReadConfig()
	collector := entity.NewCollector("http://"+cfg.Address, cfg.PollInterval, cfg.ReportInterval)

	collector.Do()

}
