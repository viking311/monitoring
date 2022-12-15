package entity

import (
	"math/rand"
	"runtime"

	"github.com/viking311/monitoring/internal/logger"
)

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
	mec.Collection["LastGC"] = &GaugeMetricEntity{Name: "LastGC", Value: float64(stat.LastGC)}
	mec.Collection["Lookups"] = &GaugeMetricEntity{Name: "Lookups", Value: float64(stat.Lookups)}
	mec.Collection["MSpanInuse"] = &GaugeMetricEntity{Name: "MSpanInuse", Value: float64(stat.MSpanInuse)}
	mec.Collection["RandomValue"] = &GaugeMetricEntity{Name: "RandomValue", Value: rand.Float64()}

	if _, ok := mec.Collection["PollCount"]; ok {
		mec.Collection["PollCount"].SetValue(uint64(1))
	} else {
		mec.Collection["PollCount"] = &CounterMetricEntity{Name: "PollCount", Value: 1}
	}
	logger.Logger.Debug("metrics updated")
}

func NewMertricCollection() MetricEntityCollection {
	return MetricEntityCollection{
		Collection: make(map[string]MetricEntityInterface),
	}

}
