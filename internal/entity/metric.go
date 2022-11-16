package entity

import (
	"fmt"
	"math/rand"
	"runtime"
)

type MetricEntityInterface interface {
	GetUpdateUri() string
}

type GuageMetricEntity struct {
	Name  string
	Value float64
}

func (gme GuageMetricEntity) GetUpdateUri() string {
	return fmt.Sprintf("/update/guage/%s/%f", gme.Name, gme.Value)
}

func (gme GuageMetricEntity) GetValue() float64 {
	return gme.Value
}

func (gme *GuageMetricEntity) SetValue(newValue float64) {
	gme.Value = newValue
}

type CounterMetricEntity struct {
	Name  string
	Value uint64
}

func (cme CounterMetricEntity) GetUpdateUri() string {
	return fmt.Sprintf("/update/counter/%s/%d", cme.Name, cme.Value)
}

type MetricEntityCollection struct {
	Collection map[string]MetricEntityInterface
	counter    uint64
}

func (mec *MetricEntityCollection) UpdateMetric(stat runtime.MemStats) {
	mec.Collection["Alloc"] = GuageMetricEntity{Name: "Alloc", Value: float64(stat.Alloc)}
	mec.Collection["BuckHashSys"] = GuageMetricEntity{Name: "BuckHashSys", Value: float64(stat.BuckHashSys)}
	mec.Collection["Frees"] = GuageMetricEntity{Name: "Frees", Value: float64(stat.Frees)}
	mec.Collection["GCCPUFraction"] = GuageMetricEntity{Name: "GCCPUFraction", Value: float64(stat.GCCPUFraction)}
	mec.Collection["GCSys"] = GuageMetricEntity{Name: "GCSys", Value: float64(stat.GCSys)}
	mec.Collection["HeapAlloc"] = GuageMetricEntity{Name: "HeapAlloc", Value: float64(stat.HeapAlloc)}
	mec.Collection["HeapIdle"] = GuageMetricEntity{Name: "HeapIdle", Value: float64(stat.HeapIdle)}
	mec.Collection["HeapInuse"] = GuageMetricEntity{Name: "HeapInuse", Value: float64(stat.HeapInuse)}
	mec.Collection["HeapObjects"] = GuageMetricEntity{Name: "HeapObjects", Value: float64(stat.HeapObjects)}
	mec.Collection["HeapReleased"] = GuageMetricEntity{Name: "HeapReleased", Value: float64(stat.HeapReleased)}
	mec.Collection["HeapSys"] = GuageMetricEntity{Name: "HeapSys", Value: float64(stat.HeapSys)}
	mec.Collection["MCacheInuse"] = GuageMetricEntity{Name: "MCacheInuse", Value: float64(stat.MCacheInuse)}
	mec.Collection["MCacheSys"] = GuageMetricEntity{Name: "MCacheSys", Value: float64(stat.MCacheSys)}
	mec.Collection["MSpanSys"] = GuageMetricEntity{Name: "MSpanSys", Value: float64(stat.MSpanSys)}
	mec.Collection["Mallocs"] = GuageMetricEntity{Name: "Mallocs", Value: float64(stat.Mallocs)}
	mec.Collection["NextGC"] = GuageMetricEntity{Name: "NextGC", Value: float64(stat.NextGC)}
	mec.Collection["NumForcedGC"] = GuageMetricEntity{Name: "NumForcedGC", Value: float64(stat.NumForcedGC)}
	mec.Collection["NumGC"] = GuageMetricEntity{Name: "NumGC", Value: float64(stat.NumGC)}
	mec.Collection["OtherSys"] = GuageMetricEntity{Name: "OtherSys", Value: float64(stat.OtherSys)}
	mec.Collection["PauseTotalNs"] = GuageMetricEntity{Name: "PauseTotalNs", Value: float64(stat.PauseTotalNs)}
	mec.Collection["StackInuse"] = GuageMetricEntity{Name: "StackInuse", Value: float64(stat.StackInuse)}
	mec.Collection["StackSys"] = GuageMetricEntity{Name: "StackSys", Value: float64(stat.StackSys)}
	mec.Collection["Sys"] = GuageMetricEntity{Name: "Sys", Value: float64(stat.Sys)}
	mec.Collection["TotalAlloc"] = GuageMetricEntity{Name: "TotalAlloc", Value: float64(stat.TotalAlloc)}
	mec.Collection["RandomValue"] = GuageMetricEntity{Name: "RandomValue", Value: rand.Float64()}

	mec.counter += 1
	mec.Collection["PollCount"] = CounterMetricEntity{Name: "PollCount", Value: mec.counter}
}

func NewMertricCollection() MetricEntityCollection {
	return MetricEntityCollection{
		Collection: make(map[string]MetricEntityInterface),
		counter:    0,
	}

}
