package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/viking311/monitoring/internal/metrics"
	"github.com/viking311/monitoring/internal/signals"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {

	mc := metrics.NewMertricCollection("http://localhost:8080/")
	sl := signals.NewSignalListener(syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	updateTicker := time.NewTicker(pollInterval)
	sendTicker := time.NewTicker(reportInterval)
	for {
		select {
		case <-updateTicker.C:
			mc.Update()
		case <-sendTicker.C:
			err := mc.SendReport()
			if err != nil {
				fmt.Println(err)
			}
		case <-sl.C:
			os.Exit(0)
		}
	}
}
