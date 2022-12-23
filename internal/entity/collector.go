package entity

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
	"github.com/viking311/monitoring/internal/logger"
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
		logger.Error(err)
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
	logger.WithFields(logrus.Fields{
		"request": request,
		"resonse": resp,
	}).Info("Metrics were sended")

	resp.Body.Close()

	return nil
}

func (c *Collector) updateStat(wg *sync.WaitGroup) {
	defer wg.Done()

	c.mtx.Lock()
	defer c.mtx.Unlock()

	runtime.ReadMemStats(&c.stat)
	c.statCollection.UpdateMetric(c.stat)
}

func (c *Collector) updateMemStat(wg *sync.WaitGroup) {
	defer wg.Done()

	c.mtx.Lock()
	defer c.mtx.Unlock()

	v, _ := mem.VirtualMemory()
	c.statCollection.UpdateMemStat(v)
}

func (c *Collector) updateCPUStat(wg *sync.WaitGroup) {
	defer wg.Done()

	c.mtx.Lock()
	defer c.mtx.Unlock()

	cpu, err := cpu.Percent(time.Millisecond, true)
	if err != nil {
		cpu = make([]float64, runtime.NumCPU())
		logger.Error(err)
	}
	c.statCollection.UpdateCPUStat(cpu)
}

// func (c *Collector) updateStatWorker(ch <-chan struct{}) {
// 	for range ch {
// 		c.updateStat()
// 	}
// }

// func (c *Collector) updateCPUMemWorker(ch <-chan struct{}) {
// 	for range ch {
// 		c.updateMemStat()
// 		c.updateCPUStat()
// 	}
// }

func (c *Collector) sendReportWorker(ch <-chan struct{}) {
	for range ch {
		c.sendReport()
	}
}

func (c *Collector) Do() {
	logger.Info("start metrics watching")
	updateTicker := time.NewTicker(c.pollInterval)
	reportTicker := time.NewTicker(c.reportInterval)

	reportCh := make(chan struct{})
	go c.sendReportWorker(reportCh)

	defer func() {
		updateTicker.Stop()
		reportTicker.Stop()
	}()

	for {
		select {
		case <-updateTicker.C:
			wg := &sync.WaitGroup{}
			wg.Add(3)

			go c.updateStat(wg)
			go c.updateMemStat(wg)
			go c.updateCPUStat(wg)

			wg.Wait()
		case <-reportTicker.C:
			reportCh <- struct{}{}
		case sig := <-c.signals.C:
			logger.WithField("signal", sig).Info("agent interrupted")
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

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go collector.updateStat(wg)
	go collector.updateMemStat(wg)
	go collector.updateCPUStat(wg)

	wg.Wait()

	return &collector
}
