package metrics

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
)

type gauge float64

type counter int64

type MetricCollection struct {
	gaugeMetric   map[string]gauge
	counterMetric map[string]counter
	endpoint      string
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

func (mc MetricCollection) SendReport() error {

	for n, v := range mc.gaugeMetric {
		vStr := fmt.Sprintf("%f", v)
		url := fmt.Sprintf("%s%s/%s", mc.endpoint, n, vStr)
		err := sendStatRequest(url, vStr)
		if err != nil {
			return err
		}
	}
	for n, v := range mc.counterMetric {
		vStr := fmt.Sprintf("%d", v)
		url := fmt.Sprintf("%s%s/%s", mc.endpoint, n, vStr)
		err := sendStatRequest(url, vStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewMertricCollection(endpoint string) MetricCollection {
	var mc = MetricCollection{
		gaugeMetric:   make(map[string]gauge),
		counterMetric: make(map[string]counter),
		endpoint:      endpoint,
	}
	mc.Update()
	return mc
}

func sendStatRequest(url string, body string) error {

	client := &http.Client{}

	bytesValue := []byte(body)
	reader := bytes.NewReader(bytesValue)
	request, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return err
	}
	fmt.Println(url)
	request.Header.Add("application-type", "text/plain")

	resp, clientErr := client.Do(request)
	if clientErr != nil {
		return clientErr
	}
	resp.Body.Close()

	return nil
}
