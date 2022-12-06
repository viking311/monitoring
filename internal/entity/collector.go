package entity

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	hashKey        string
}

func (c *Collector) sendReport() {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for _, metric := range c.statCollection.Collection {
		go c.sendStatRequest(metric.GetUpdateURI(), metric.GetMetricEntity())
	}
	c.statCollection.Collection["PollCount"] = &CounterMetricEntity{Name: "PollCount", Value: 0}
}

func (c *Collector) sendStatRequest(uri string, value Metrics) {
	c.calculateHash(&value)
	bytesValue, err := json.Marshal(value)
	if err != nil {
		return
	}

	reader := bytes.NewReader(bytesValue)
	request, err := http.NewRequest(http.MethodPost, c.endpoint+uri, reader)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

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

func (c *Collector) calculateHash(data *Metrics) {
	if len(c.hashKey) > 0 {
		src := ""
		hasher := hmac.New(sha256.New, []byte(c.hashKey))
		if data.MType == "counter" {
			src = fmt.Sprintf("%s:counter:%d", data.ID, data.Delta)
		}

		if data.MType == "gauge" {
			src = fmt.Sprintf("%s:gauge:%f", data.ID, *data.Value)
		}
		if len(src) > 0 {
			hasher.Write([]byte(src))
			data.Hash = hex.EncodeToString(hasher.Sum(nil))
		}

	}
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

func NewCollector(endpoint string, pollInterval time.Duration, reportInterval time.Duration, hashKey string) *Collector {
	var collector = Collector{
		endpoint:       strings.TrimSuffix(endpoint, "/"),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		statCollection: NewMertricCollection(),
		stat:           runtime.MemStats{},
		signals:        signals.NewSignalListener(syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT),
		mtx:            sync.RWMutex{},
		hashKey:        hashKey,
	}

	collector.updateStat()

	return &collector
}
