package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/logger"
	"github.com/viking311/monitoring/internal/signals"
	"github.com/viking311/monitoring/internal/storage"
)

type Collector struct {
	endpoint       string
	pollInterval   time.Duration
	reportInterval time.Duration
	storage        storage.Repository
	signals        signals.SignalListener
	mtx            sync.RWMutex
	hashKey        string
}

func (c *Collector) Do() {
	logger.Info("start metrics watching")

	updateTicker := time.NewTicker(c.pollInterval)
	reportTicker := time.NewTicker(c.reportInterval)
	defer func() {
		updateTicker.Stop()
		reportTicker.Stop()
	}()

	reportCh := make(chan struct{})
	go c.sendReportWorker(reportCh)

	for {
		select {
		case <-updateTicker.C:
			c.updateMetrics()
		case <-reportTicker.C:
			reportCh <- struct{}{}
		case sig := <-c.signals.C:
			logger.WithField("signal", sig).Info("agent interrupted")
			return
		}
	}

}

func (c *Collector) updateMetrics() {
	logger.Info("update metrics")

	c.mtx.Lock()
	defer c.mtx.Unlock()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go c.updateMemStat(wg)
	go c.updateVirtualMemoryStat(wg)
	go c.updateCPUStat(wg)

	wg.Wait()
}

func (c *Collector) updateMemStat(wg *sync.WaitGroup) {
	logger.Debug("update memStat")

	defer wg.Done()

	stat := runtime.MemStats{}
	runtime.ReadMemStats(&stat)
	values := make([]entity.Metrics, 29)
	alloc := float64(stat.Alloc)
	value := entity.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &alloc,
	}
	values = append(values, value)

	buckHashSys := float64(stat.BuckHashSys)
	value = entity.Metrics{
		ID:    "BuckHashSys",
		MType: "gauge",
		Value: &buckHashSys,
	}
	values = append(values, value)

	frees := float64(stat.Frees)
	value = entity.Metrics{
		ID:    "Frees",
		MType: "gauge",
		Value: &frees,
	}
	values = append(values, value)

	gCCPUFraction := float64(stat.GCCPUFraction)
	value = entity.Metrics{
		ID:    "GCCPUFraction",
		MType: "gauge",
		Value: &gCCPUFraction,
	}
	values = append(values, value)

	gcSys := float64(stat.GCSys)
	value = entity.Metrics{
		ID:    "GCSys",
		MType: "gauge",
		Value: &gcSys,
	}
	values = append(values, value)

	heapAlloc := float64(stat.HeapAlloc)
	value = entity.Metrics{
		ID:    "HeapAlloc",
		MType: "gauge",
		Value: &heapAlloc,
	}
	values = append(values, value)

	heapIdle := float64(stat.HeapIdle)
	value = entity.Metrics{
		ID:    "HeapIdle",
		MType: "gauge",
		Value: &heapIdle,
	}
	values = append(values, value)

	heapInuse := float64(stat.HeapInuse)
	value = entity.Metrics{
		ID:    "HeapInuse",
		MType: "gauge",
		Value: &heapInuse,
	}
	values = append(values, value)

	heapObjects := float64(stat.HeapObjects)
	value = entity.Metrics{
		ID:    "HeapObjects",
		MType: "gauge",
		Value: &heapObjects,
	}
	values = append(values, value)

	heapReleased := float64(stat.HeapReleased)
	value = entity.Metrics{
		ID:    "HeapReleased",
		MType: "gauge",
		Value: &heapReleased,
	}
	values = append(values, value)

	heapSys := float64(stat.HeapSys)
	value = entity.Metrics{
		ID:    "HeapSys",
		MType: "gauge",
		Value: &heapSys,
	}
	values = append(values, value)

	mCacheInuse := float64(stat.MCacheInuse)
	value = entity.Metrics{
		ID:    "MCacheInuse",
		MType: "gauge",
		Value: &mCacheInuse,
	}
	values = append(values, value)

	mCacheSys := float64(stat.MCacheSys)
	value = entity.Metrics{
		ID:    "MCacheSys",
		MType: "gauge",
		Value: &mCacheSys,
	}
	values = append(values, value)

	mSpanSys := float64(stat.MSpanSys)
	value = entity.Metrics{
		ID:    "MSpanSys",
		MType: "gauge",
		Value: &mSpanSys,
	}
	values = append(values, value)

	mallocs := float64(stat.Mallocs)
	value = entity.Metrics{
		ID:    "Mallocs",
		MType: "gauge",
		Value: &mallocs,
	}
	values = append(values, value)

	nextGC := float64(stat.NextGC)
	value = entity.Metrics{
		ID:    "NextGC",
		MType: "gauge",
		Value: &nextGC,
	}
	values = append(values, value)

	numForcedGC := float64(stat.NumForcedGC)
	value = entity.Metrics{
		ID:    "NumForcedGC",
		MType: "gauge",
		Value: &numForcedGC,
	}
	values = append(values, value)

	numGC := float64(stat.NumGC)
	value = entity.Metrics{
		ID:    "NumGC",
		MType: "gauge",
		Value: &numGC,
	}
	values = append(values, value)

	otherSys := float64(stat.OtherSys)
	value = entity.Metrics{
		ID:    "OtherSys",
		MType: "gauge",
		Value: &otherSys,
	}
	values = append(values, value)

	pauseTotalNs := float64(stat.PauseTotalNs)
	value = entity.Metrics{
		ID:    "PauseTotalNs",
		MType: "gauge",
		Value: &pauseTotalNs,
	}
	values = append(values, value)

	stackInuse := float64(stat.StackInuse)
	value = entity.Metrics{
		ID:    "StackInuse",
		MType: "gauge",
		Value: &stackInuse,
	}
	values = append(values, value)

	stackSys := float64(stat.StackSys)
	value = entity.Metrics{
		ID:    "StackSys",
		MType: "gauge",
		Value: &stackSys,
	}
	values = append(values, value)

	sys := float64(stat.Sys)
	value = entity.Metrics{
		ID:    "Sys",
		MType: "gauge",
		Value: &sys,
	}
	values = append(values, value)

	totalAlloc := float64(stat.TotalAlloc)
	value = entity.Metrics{
		ID:    "TotalAlloc",
		MType: "gauge",
		Value: &totalAlloc,
	}
	values = append(values, value)

	lastGC := float64(stat.LastGC)
	value = entity.Metrics{
		ID:    "LastGC",
		MType: "gauge",
		Value: &lastGC,
	}
	values = append(values, value)

	lookups := float64(stat.Lookups)
	value = entity.Metrics{
		ID:    "Lookups",
		MType: "gauge",
		Value: &lookups,
	}
	values = append(values, value)

	mSpanInuse := float64(stat.MSpanInuse)
	value = entity.Metrics{
		ID:    "MSpanInuse",
		MType: "gauge",
		Value: &mSpanInuse,
	}
	values = append(values, value)

	randomValue := rand.Float64()
	value = entity.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randomValue,
	}
	values = append(values, value)

	pollCount := uint64(1)
	value = entity.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	values = append(values, value)

	err := c.storage.BatchUpdate(values)
	if err != nil {
		logger.Error(err)
	}
}

func (c *Collector) updateVirtualMemoryStat(wg *sync.WaitGroup) {
	logger.Debug("update VirtualMemoryStat")
	defer wg.Done()
	v, err := mem.VirtualMemory()
	if err != nil {
		logger.Error(err)
		return
	}
	values := make([]entity.Metrics, 2)

	totalMemory := float64(v.Total)
	value := entity.Metrics{
		ID:    "TotalMemory",
		MType: "gauge",
		Value: &totalMemory,
	}
	values = append(values, value)

	free := float64(v.Free)
	value = entity.Metrics{
		ID:    "FreeMemory",
		MType: "gauge",
		Value: &free,
	}
	values = append(values, value)
	err = c.storage.BatchUpdate(values)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (c *Collector) updateCPUStat(wg *sync.WaitGroup) {
	logger.Debug("update CPUStat")

	defer wg.Done()

	cpu, err := cpu.Percent(time.Millisecond, true)
	if err != nil {
		logger.Error(err)
		return
	}

	values := make([]entity.Metrics, runtime.NumCPU())

	for ind, val := range cpu {
		value := entity.Metrics{
			ID:    fmt.Sprintf("CPUutilization%d", ind+1),
			MType: "gauge",
			Value: &val,
		}
		values = append(values, value)
	}
	fmt.Println(values)
	// for i := 0; i < runtime.NumCPU(); i++ {
	// 	value := entity.Metrics{
	// 		ID:    fmt.Sprintf("CPUutilization%d", i+1),
	// 		MType: "gauge",
	// 		Value: &cpu[i],
	// 	}
	// 	values = append(values, value)
	// }

	err = c.storage.BatchUpdate(values)
	if err != nil {
		logger.Error(err)
		return
	}

}

func (c *Collector) sendReportWorker(ch <-chan struct{}) {
	for range ch {
		c.sendReport()
	}
}

func (c *Collector) sendReport() {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	collection, err := c.storage.GetAll()
	if err != nil {
		logger.Error(err)
		return
	}

	if len(collection) == 0 {
		return
	}

	values := make([]entity.Metrics, 0, len(collection))
	for _, item := range collection {
		item.Hash = entity.MetricsHash(item, c.hashKey)
		values = append(values, item)
	}

	err = c.sendBatchRequest(values)
	if err != nil {
		logger.Error(err)
	}
	value := entity.Metrics{
		ID:    "PollCount",
		MType: "counter",
	}

	err = c.storage.Delete(value.GetKey())
	if err != nil {
		logger.Error(err)
	}

}

func (c *Collector) sendBatchRequest(values []entity.Metrics) error {
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
	}).Info("Metrics were sent")

	resp.Body.Close()

	return nil
}

func NewCollector(endpoint string, pollInterval time.Duration, reportInterval time.Duration, hashKey string) *Collector {
	var collector = Collector{
		endpoint:       strings.TrimSuffix(endpoint, "/"),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		storage:        storage.NewInMemoryStorage(false),
		signals:        signals.NewSignalListener(syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT),
		mtx:            sync.RWMutex{},
		hashKey:        hashKey,
	}

	collector.updateMetrics()

	return &collector

}
