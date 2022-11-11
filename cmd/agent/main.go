package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type gauge float64

type counter int64

type MetricCollection struct {
	gaugeMetric   map[string]gauge
	counterMetric map[string]counter
}

func (mc *MetricCollection) Update() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mc.gaugeMetric["Alloc"] = gauge(m.Alloc)
	mc.gaugeMetric["BuckHashSys"] = gauge(m.BuckHashSys)
	mc.gaugeMetric["Frees"] = gauge(m.Frees)
	mc.gaugeMetric["GCCPUFraction"] = gauge(m.GCCPUFraction)
	mc.gaugeMetric["GCSys"] = gauge(m.GCSys)
	mc.gaugeMetric["HeapAlloc"] = gauge(m.HeapAlloc)
	mc.gaugeMetric["HeapIdle"] = gauge(m.HeapIdle)
	mc.gaugeMetric["HeapInuse"] = gauge(m.HeapInuse)
	mc.gaugeMetric["HeapObjects"] = gauge(m.HeapObjects)
	mc.gaugeMetric["HeapReleased"] = gauge(m.HeapReleased)
	mc.gaugeMetric["HeapSys"] = gauge(m.HeapSys)
	mc.gaugeMetric["MCacheInuse"] = gauge(m.MCacheInuse)
	mc.gaugeMetric["MCacheSys"] = gauge(m.MCacheSys)
	mc.gaugeMetric["MSpanSys"] = gauge(m.MSpanSys)
	mc.gaugeMetric["Mallocs"] = gauge(m.Mallocs)
	mc.gaugeMetric["NextGC"] = gauge(m.NextGC)
	mc.gaugeMetric["NumForcedGC"] = gauge(m.NumForcedGC)
	mc.gaugeMetric["NumGC"] = gauge(m.NumGC)
	mc.gaugeMetric["OtherSys"] = gauge(m.OtherSys)
	mc.gaugeMetric["PauseTotalNs"] = gauge(m.PauseTotalNs)
	mc.gaugeMetric["StackInuse"] = gauge(m.StackInuse)
	mc.gaugeMetric["StackSys"] = gauge(m.StackSys)
	mc.gaugeMetric["Sys"] = gauge(m.Sys)
	mc.gaugeMetric["TotalAlloc"] = gauge(m.TotalAlloc)
	mc.gaugeMetric["RandomValue"] = gauge(rand.Float64())
	mc.counterMetric["PollCount"] += 1
}

func (mc MetricCollection) SendReport() {
	var requests [100]http.Request
	i := 0
	client := &http.Client{}
	for n, v := range mc.gaugeMetric {
		vStr := fmt.Sprintf("%f", v)
		bytesValue := []byte(vStr)
		reader := bytes.NewReader(bytesValue)
		url := fmt.Sprintf("%s%s/%s", endpoint, n, vStr)
		request, err := http.NewRequest(http.MethodPost, url, reader)
		request.Header.Add("application-type", "text/plain")
		if err != nil {
			continue
		}
		// fmt.Print("Send report to url: " + url)
		_, clientErr := client.Do(request)
		if clientErr != nil {
			// fmt.Println(err)
			continue
		}
		// fmt.Println(" ResultCode:", res.Status)
		requests[i] = *request
		i++

	}
	for n, v := range mc.counterMetric {
		vStr := fmt.Sprintf("%d", v)
		bytesValue := []byte(vStr)
		reader := bytes.NewReader(bytesValue)
		url := fmt.Sprintf("%s%s/%s", endpoint, n, vStr)
		request, err := http.NewRequest(http.MethodPost, url, reader)
		request.Header.Add("application-type", "text/plain")
		if err != nil {
			continue
		}
		// fmt.Print("Send report to url: " + url)
		_, clientErr := client.Do(request)
		if clientErr != nil {
			// fmt.Println(err)
			continue
		}
		// fmt.Println(" ResultCode:", res.Status)
		// fmt.Println("Send report to url: " + fmt.Sprintf("%s%s/%d", endpoint, n, v))
	}

}

const (
	endpoint       = "http://localhost:8080/"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {

	var mc = MetricCollection{
		gaugeMetric:   make(map[string]gauge),
		counterMetric: make(map[string]counter),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	mc.Update()
	fmt.Println(mc)
	updateTicker := time.NewTicker(pollInterval)
	sendTicker := time.NewTicker(reportInterval)
	for {
		select {
		case <-updateTicker.C:
			mc.Update()
			fmt.Println(mc)
		case <-sendTicker.C:
			mc.SendReport()
		case <-c:
			os.Exit(0)
		}
	}
}
