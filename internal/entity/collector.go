package entity

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
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
	if len(c.statCollection.Collection) == 0 {
		return
	}
	values := make([]Metrics, 0, len(c.statCollection.Collection))

	for _, metric := range c.statCollection.Collection {
		mc := metric.GetMetricEntity()
		mc.Hash = MetricsHash(mc, c.hashKey)
		values = append(values, mc)
	}
	err := c.sendBatchRequest(values)
	if err != nil {
		log.Println(err)
	}
	c.statCollection.Collection["PollCount"] = &CounterMetricEntity{Name: "PollCount", Value: 0}
}

func (c *Collector) sendBatchRequest(values []Metrics) error {
	bytesValue, err := json.Marshal(values)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	gzWriter, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return err
	}
	_, err = gzWriter.Write(bytesValue)
	if err != nil {
		return err
	}
	gzWriter.Close()

	reader := bytes.NewReader(b.Bytes())
	request, err := http.NewRequest(http.MethodPost, c.endpoint+"/updates/", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{}
	resp, clientErr := client.Do(request)
	if clientErr != nil {
		return clientErr
	}
	log.Println(resp)

	resp.Body.Close()
	return nil
}

func (c *Collector) updateStat() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	runtime.ReadMemStats(&c.stat)
	c.statCollection.UpdateMetric(c.stat)
}

func (c *Collector) Do() {
	log.Println("start metrics watching")
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
			log.Println("agent interrupted")
			return
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
