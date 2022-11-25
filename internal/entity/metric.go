package entity

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
)

type MetricEntityInterface interface {
	GetUpdateURI() string
	GetValue() interface{}
	SetValue(value interface{})
	GetKey() string
	GetStringValue() string
	GetShortTypeName() string
}

type GaugeMetricEntity struct {
	Name  string
	Value float64
}

func (gme *GaugeMetricEntity) GetUpdateURI() string {
	return fmt.Sprintf("/update/gauge/%s/%f", gme.Name, gme.Value)
}

func (gme *GaugeMetricEntity) GetValue() interface{} {
	return gme.Value
}

func (gme *GaugeMetricEntity) GetKey() string {
	return gme.Name
}

func (gme *GaugeMetricEntity) SetValue(value interface{}) {
	floatValue, ok := value.(float64)
	if ok {
		gme.Value = floatValue
	}
}

func (gme *GaugeMetricEntity) GetStringValue() string {
	return strconv.FormatFloat(gme.Value, 'f', -1, 64)
}

func (gme *GaugeMetricEntity) GetShortTypeName() string {
	return "gauge"
}

type CounterMetricEntity struct {
	Name  string
	Value uint64
}

func (cme *CounterMetricEntity) GetUpdateURI() string {
	return fmt.Sprintf("/update/counter/%s/%d", cme.Name, cme.Value)
}

func (cme *CounterMetricEntity) GetValue() interface{} {
	return cme.Value
}

func (cme *CounterMetricEntity) SetValue(value interface{}) {
	intValue, ok := value.(uint64)
	if ok {
		cme.Value += intValue
	}
}

func (cme *CounterMetricEntity) GetKey() string {
	return cme.Name
}

func (cme *CounterMetricEntity) GetStringValue() string {
	return fmt.Sprintf("%d", cme.Value)
}

func (cme *CounterMetricEntity) GetShortTypeName() string {
	return "counter"
}

type MetricEntityCollection struct {
	Collection map[string]MetricEntityInterface
}

func (mec *MetricEntityCollection) UpdateMetric(stat runtime.MemStats) {
	mec.Collection["Alloc"] = &GaugeMetricEntity{Name: "Alloc", Value: float64(stat.Alloc)}
	mec.Collection["BuckHashSys"] = &GaugeMetricEntity{Name: "BuckHashSys", Value: float64(stat.BuckHashSys)}
	mec.Collection["Frees"] = &GaugeMetricEntity{Name: "Frees", Value: float64(stat.Frees)}
	mec.Collection["GCCPUFraction"] = &GaugeMetricEntity{Name: "GCCPUFraction", Value: float64(stat.GCCPUFraction)}
	mec.Collection["GCSys"] = &GaugeMetricEntity{Name: "GCSys", Value: float64(stat.GCSys)}
	mec.Collection["HeapAlloc"] = &GaugeMetricEntity{Name: "HeapAlloc", Value: float64(stat.HeapAlloc)}
	mec.Collection["HeapIdle"] = &GaugeMetricEntity{Name: "HeapIdle", Value: float64(stat.HeapIdle)}
	mec.Collection["HeapInuse"] = &GaugeMetricEntity{Name: "HeapInuse", Value: float64(stat.HeapInuse)}
	mec.Collection["HeapObjects"] = &GaugeMetricEntity{Name: "HeapObjects", Value: float64(stat.HeapObjects)}
	mec.Collection["HeapReleased"] = &GaugeMetricEntity{Name: "HeapReleased", Value: float64(stat.HeapReleased)}
	mec.Collection["HeapSys"] = &GaugeMetricEntity{Name: "HeapSys", Value: float64(stat.HeapSys)}
	mec.Collection["MCacheInuse"] = &GaugeMetricEntity{Name: "MCacheInuse", Value: float64(stat.MCacheInuse)}
	mec.Collection["MCacheSys"] = &GaugeMetricEntity{Name: "MCacheSys", Value: float64(stat.MCacheSys)}
	mec.Collection["MSpanSys"] = &GaugeMetricEntity{Name: "MSpanSys", Value: float64(stat.MSpanSys)}
	mec.Collection["Mallocs"] = &GaugeMetricEntity{Name: "Mallocs", Value: float64(stat.Mallocs)}
	mec.Collection["NextGC"] = &GaugeMetricEntity{Name: "NextGC", Value: float64(stat.NextGC)}
	mec.Collection["NumForcedGC"] = &GaugeMetricEntity{Name: "NumForcedGC", Value: float64(stat.NumForcedGC)}
	mec.Collection["NumGC"] = &GaugeMetricEntity{Name: "NumGC", Value: float64(stat.NumGC)}
	mec.Collection["OtherSys"] = &GaugeMetricEntity{Name: "OtherSys", Value: float64(stat.OtherSys)}
	mec.Collection["PauseTotalNs"] = &GaugeMetricEntity{Name: "PauseTotalNs", Value: float64(stat.PauseTotalNs)}
	mec.Collection["StackInuse"] = &GaugeMetricEntity{Name: "StackInuse", Value: float64(stat.StackInuse)}
	mec.Collection["StackSys"] = &GaugeMetricEntity{Name: "StackSys", Value: float64(stat.StackSys)}
	mec.Collection["Sys"] = &GaugeMetricEntity{Name: "Sys", Value: float64(stat.Sys)}
	mec.Collection["TotalAlloc"] = &GaugeMetricEntity{Name: "TotalAlloc", Value: float64(stat.TotalAlloc)}
	mec.Collection["RandomValue"] = &GaugeMetricEntity{Name: "RandomValue", Value: rand.Float64()}

	if _, ok := mec.Collection["PollCount"]; ok {
		mec.Collection["PollCount"].SetValue(uint64(1))
	} else {
		mec.Collection["PollCount"] = &CounterMetricEntity{Name: "PollCount", Value: 1}
	}

}

func NewMertricCollection() MetricEntityCollection {
	return MetricEntityCollection{
		Collection: make(map[string]MetricEntityInterface),
	}

}
