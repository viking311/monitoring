package entity

import (
	"bytes"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
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
	mtx            sync.RWMutex
}

func (c *Collector) sendReport() {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for _, metric := range c.statCollection.Collection {
		go c.sendStatRequest(metric.GetUpdateURI(), metric.GetStringValue())
	}
	c.statCollection.Collection["PollCount"] = &CounterMetricEntity{Name: "PollCount", Value: 0}
}

func (c *Collector) sendStatRequest(uri string, value string) {
	bytesValue := []byte(value)
	reader := bytes.NewReader(bytesValue)
	request, err := http.NewRequest(http.MethodPost, c.endpoint+uri, reader)
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
	c.mtx.Lock()
	defer c.mtx.Unlock()

	runtime.ReadMemStats(&c.stat)
	c.statCollection.UpdateMetric(c.stat)
}

func (c *Collector) Do() {
	updateTicker := time.NewTicker(c.pollInterval)
	reportTicker := time.NewTicker(c.reportInterval)

	defer func() {
		updateTicker.Stop()
		reportTicker.Stop()
	}()

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

func NewCollector(endpoint string, pollInterval time.Duration, reportInterval time.Duration) *Collector {
	var collector = Collector{
		endpoint:       strings.TrimSuffix(endpoint, "/"),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		statCollection: NewMertricCollection(),
		stat:           runtime.MemStats{},
		signals:        signals.NewSignalListener(syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT),
		mtx:            sync.RWMutex{},
	}

	collector.updateStat()

	return &collector
}
