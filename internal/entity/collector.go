package entity

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/viking311/monitoring/internal/signals"
)

type Collector struct {
	endpoint       string
	pollInterval   time.Duration
	reportInterval time.Duration
	statCollection MetricEntityCollection
	stat           runtime.MemStats
	signals        signals.SignalListener
}

func (c *Collector) sendReport() {
	for _, metric := range c.statCollection.Collection {
		go c.sendStatRequest(metric.GetUpdateURI(), metric.GetStringValue())
	}
}

func (c *Collector) sendStatRequest(uri string, value string) {
	fmt.Println(uri)
	bytesValue := []byte(value)
	reader := bytes.NewReader(bytesValue)
	request, err := http.NewRequest(http.MethodPost, c.endpoint+"/"+uri, reader)
	if err != nil {
		return
	}
	request.Header.Set("application-type", "text/plain")

	client := &http.Client{}
	resp, clientErr := client.Do(request)
	if clientErr != nil {
		return
	}
	defer resp.Body.Close()
}

func (c *Collector) updateStat() {
	runtime.ReadMemStats(&c.stat)
	c.statCollection.UpdateMetric(c.stat)
}

func (c *Collector) Do() {
	updateTicker := time.NewTicker(c.pollInterval)
	reportTicker := time.NewTicker(c.reportInterval)
	for {
		select {
		case <-updateTicker.C:
			c.updateStat()
		case <-reportTicker.C:
			c.sendReport()
		case <-c.signals.C:
			os.Exit(0)
		}

	}

}

func NewCollector(endpoint string, pollInterval time.Duration, reportInterval time.Duration) Collector {
	var collector = Collector{
		endpoint:       strings.TrimSuffix(endpoint, "/"),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		statCollection: NewMertricCollection(),
		stat:           runtime.MemStats{},
		signals:        signals.NewSignalListener(syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT),
	}

	collector.updateStat()

	return collector
}
