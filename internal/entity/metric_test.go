package entity

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricEntityCollection_UpdateMetric(t *testing.T) {
	type args struct {
		stat runtime.MemStats
	}
	mec := NewMertricCollection()

	stat := runtime.MemStats{}

	stat.Alloc = 1
	stat.BuckHashSys = 1
	stat.Frees = 1
	stat.GCCPUFraction = 1
	stat.GCSys = 1
	stat.HeapAlloc = 1
	stat.HeapIdle = 1
	stat.HeapInuse = 1
	stat.HeapObjects = 1
	stat.HeapReleased = 1
	stat.HeapSys = 1
	stat.MCacheInuse = 1
	stat.MCacheSys = 1
	stat.MSpanSys = 1
	stat.Mallocs = 1
	stat.NextGC = 1
	stat.NumForcedGC = 1
	stat.NumGC = 1
	stat.OtherSys = 1
	stat.PauseTotalNs = 1
	stat.StackInuse = 1
	stat.StackSys = 1
	stat.Sys = 1
	stat.TotalAlloc = 1
	stat.LastGC = 1
	stat.Lookups = 1
	stat.MSpanInuse = 1

	tests := []struct {
		name string
		mec  *MetricEntityCollection
		args args
		want *MetricEntityCollection
	}{
		{
			name: "UpdateMetric: test1",
			mec:  &mec,
			args: args{
				stat: stat,
			},
			want: getUpdateResult(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mec.UpdateMetric(tt.args.stat)
			assert.Equal(t, 29, len(tt.mec.Collection), "too few metrics")
			fmt.Println(tt.mec.Collection)
			_, ok := tt.mec.Collection["RandomValue"]
			assert.Equal(t, true, ok)
			delete(tt.mec.Collection, "RandomValue")
			assert.Equal(t, tt.want, tt.mec)
		})
	}
}
func getUpdateResult() *MetricEntityCollection {
	mce := MetricEntityCollection{
		Collection: make(map[string]MetricEntityInterface),
	}

	mce.Collection["Alloc"] = &GaugeMetricEntity{Name: "Alloc", Value: float64(1)}
	mce.Collection["BuckHashSys"] = &GaugeMetricEntity{Name: "BuckHashSys", Value: float64(1)}
	mce.Collection["Frees"] = &GaugeMetricEntity{Name: "Frees", Value: float64(1)}
	mce.Collection["GCCPUFraction"] = &GaugeMetricEntity{Name: "GCCPUFraction", Value: float64(1)}
	mce.Collection["GCSys"] = &GaugeMetricEntity{Name: "GCSys", Value: float64(1)}
	mce.Collection["HeapAlloc"] = &GaugeMetricEntity{Name: "HeapAlloc", Value: float64(1)}
	mce.Collection["HeapIdle"] = &GaugeMetricEntity{Name: "HeapIdle", Value: float64(1)}
	mce.Collection["HeapInuse"] = &GaugeMetricEntity{Name: "HeapInuse", Value: float64(1)}
	mce.Collection["HeapObjects"] = &GaugeMetricEntity{Name: "HeapObjects", Value: float64(1)}
	mce.Collection["HeapReleased"] = &GaugeMetricEntity{Name: "HeapReleased", Value: float64(1)}
	mce.Collection["HeapSys"] = &GaugeMetricEntity{Name: "HeapSys", Value: float64(1)}
	mce.Collection["MCacheInuse"] = &GaugeMetricEntity{Name: "MCacheInuse", Value: float64(1)}
	mce.Collection["MCacheSys"] = &GaugeMetricEntity{Name: "MCacheSys", Value: float64(1)}
	mce.Collection["MSpanSys"] = &GaugeMetricEntity{Name: "MSpanSys", Value: float64(1)}
	mce.Collection["Mallocs"] = &GaugeMetricEntity{Name: "Mallocs", Value: float64(1)}
	mce.Collection["NextGC"] = &GaugeMetricEntity{Name: "NextGC", Value: float64(1)}
	mce.Collection["NumForcedGC"] = &GaugeMetricEntity{Name: "NumForcedGC", Value: float64(1)}
	mce.Collection["NumGC"] = &GaugeMetricEntity{Name: "NumGC", Value: float64(1)}
	mce.Collection["OtherSys"] = &GaugeMetricEntity{Name: "OtherSys", Value: float64(1)}
	mce.Collection["PauseTotalNs"] = &GaugeMetricEntity{Name: "PauseTotalNs", Value: float64(1)}
	mce.Collection["StackInuse"] = &GaugeMetricEntity{Name: "StackInuse", Value: float64(1)}
	mce.Collection["StackSys"] = &GaugeMetricEntity{Name: "StackSys", Value: float64(1)}
	mce.Collection["Sys"] = &GaugeMetricEntity{Name: "Sys", Value: float64(1)}
	mce.Collection["TotalAlloc"] = &GaugeMetricEntity{Name: "TotalAlloc", Value: float64(1)}
	mce.Collection["LastGC"] = &GaugeMetricEntity{Name: "LastGC", Value: float64(1)}
	mce.Collection["Lookups"] = &GaugeMetricEntity{Name: "Lookups", Value: float64(1)}
	mce.Collection["MSpanInuse"] = &GaugeMetricEntity{Name: "MSpanInuse", Value: float64(1)}
	mce.Collection["PollCount"] = &CounterMetricEntity{Name: "PollCount", Value: uint64(1)}

	return &mce
}

func TestMetricEntityCollection_UpdateMetric_counterUpdate(t *testing.T) {
	mec := NewMertricCollection()

	stat := runtime.MemStats{}

	stat.Alloc = 1
	stat.BuckHashSys = 1
	stat.Frees = 1
	stat.GCCPUFraction = 1
	stat.GCSys = 1
	stat.HeapAlloc = 1
	stat.HeapIdle = 1
	stat.HeapInuse = 1
	stat.HeapObjects = 1
	stat.HeapReleased = 1
	stat.HeapSys = 1
	stat.MCacheInuse = 1
	stat.MCacheSys = 1
	stat.MSpanSys = 1
	stat.Mallocs = 1
	stat.NextGC = 1
	stat.NumForcedGC = 1
	stat.NumGC = 1
	stat.OtherSys = 1
	stat.PauseTotalNs = 1
	stat.StackInuse = 1
	stat.StackSys = 1
	stat.Sys = 1
	stat.TotalAlloc = 1
	stat.LastGC = 1
	stat.Lookups = 1
	stat.MSpanInuse = 1

	mec.UpdateMetric(stat)
	mec.UpdateMetric(stat)
	assert.Equal(t, uint64(2), mec.Collection["PollCount"].GetValue().(uint64))
}

func TestGaugeMetricEntity_GetUpdateURI(t *testing.T) {
	tests := []struct {
		name string
		gme  *GaugeMetricEntity
		want string
	}{
		{
			name: "GetUpdateURI: test1",
			gme: &GaugeMetricEntity{
				Name:  "metric1",
				Value: float64(1),
			},
			// want: "/update/gauge/metric1/1.000000",
			want: "/update/",
		},
		{
			name: "GetUpdateURI: test2",
			gme: &GaugeMetricEntity{
				Name:  "metric2",
				Value: float64(1.5),
			},
			// want: "/update/gauge/metric2/1.500000",
			want: "/update/",
		},
		{
			name: "GetUpdateURI: test3",
			gme: &GaugeMetricEntity{
				Name:  "metric2",
				Value: 0,
			},
			// want: "/update/gauge/metric2/0.000000",
			want: "/update/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gme.GetUpdateURI(); got != tt.want {
				t.Errorf("GaugeMetricEntity.GetUpdateURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaugeMetricEntity_GetValue(t *testing.T) {
	tests := []struct {
		name string
		gme  *GaugeMetricEntity
		want interface{}
	}{
		{
			name: "GetValue: test1",
			gme: &GaugeMetricEntity{
				Name:  "Metric1",
				Value: 1,
			},
			want: float64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gme.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GaugeMetricEntity.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaugeMetricEntity_GetKey(t *testing.T) {
	tests := []struct {
		name string
		gme  *GaugeMetricEntity
		want string
	}{
		{
			name: "GetKey: test1",
			gme: &GaugeMetricEntity{
				Name:  "guageMetric",
				Value: 1,
			},
			want: "guagemetric",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gme.GetKey(); got != tt.want {
				t.Errorf("GaugeMetricEntity.GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaugeMetricEntity_SetValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		gme  *GaugeMetricEntity
		args args
		want float64
	}{
		{
			name: "SetValue: test1",
			gme: &GaugeMetricEntity{
				Name:  "m1",
				Value: 1,
			},
			args: args{
				value: float64(2),
			},
			want: float64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gme.SetValue(tt.args.value)
			assert.Equal(t, tt.want, tt.gme.Value)
		})
	}
}

func TestGaugeMetricEntity_GetStringValue(t *testing.T) {
	tests := []struct {
		name string
		gme  *GaugeMetricEntity
		want string
	}{
		{
			name: "TestGaugeMetricEntity_GetStringValue: test1",
			gme: &GaugeMetricEntity{
				Name:  "m1",
				Value: 1,
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gme.GetStringValue(); got != tt.want {
				t.Errorf("GaugeMetricEntity.GetStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetricEntity_GetUpdateURI(t *testing.T) {
	tests := []struct {
		name string
		cme  *CounterMetricEntity
		want string
	}{
		{
			name: "TestGaugeMetricEntity_GetStringValue: test1",
			cme: &CounterMetricEntity{
				Name:  "m1",
				Value: uint64(1),
			},
			// want: "/update/counter/m1/1",
			want: "/update/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cme.GetUpdateURI(); got != tt.want {
				t.Errorf("CounterMetricEntity.GetUpdateURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetricEntity_GetValue(t *testing.T) {
	tests := []struct {
		name string
		cme  *CounterMetricEntity
		want interface{}
	}{
		{
			name: "CounterMetricEntity_GetValue: test1",
			cme: &CounterMetricEntity{
				Name:  "m1",
				Value: uint64(1),
			},
			want: uint64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cme.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CounterMetricEntity.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetricEntity_SetValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		cme  *CounterMetricEntity
		args args
		want uint64
	}{
		{
			name: "CounterMetricEntity_SetValue: test1",
			cme: &CounterMetricEntity{
				Name:  "m1",
				Value: uint64(1),
			},
			args: args{
				value: uint64(2),
			},
			want: uint64(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cme.SetValue(tt.args.value)
			assert.Equal(t, tt.want, tt.cme.Value)
		})
	}
}

func TestCounterMetricEntity_GetKey(t *testing.T) {
	tests := []struct {
		name string
		cme  *CounterMetricEntity
		want string
	}{
		{
			name: "TestCounterMetricEntity_GetKey: test1",
			cme: &CounterMetricEntity{
				Name:  "m1",
				Value: 1,
			},
			want: "m1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cme.GetKey(); got != tt.want {
				t.Errorf("CounterMetricEntity.GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetricEntity_GetStringValue(t *testing.T) {
	tests := []struct {
		name string
		cme  *CounterMetricEntity
		want string
	}{
		{
			name: "TestCounterMetricEntity_GetStringValue: test1",
			cme: &CounterMetricEntity{
				Name:  "m1",
				Value: 1,
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cme.GetStringValue(); got != tt.want {
				t.Errorf("CounterMetricEntity.GetStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
