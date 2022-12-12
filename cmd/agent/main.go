package main

import (
	"log"

	"github.com/viking311/monitoring/internal/agent"
	"github.com/viking311/monitoring/internal/entity"
)

func main() {
	log.Println("start agent")
	defer log.Println("stop agent")

	collector := entity.NewCollector("http://"+*agent.Config.Address, *agent.Config.PollInterval, *agent.Config.ReportInterval, *agent.Config.HashKey)

	collector.Do()
}
