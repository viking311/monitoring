package entity

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollector_sendReport(t *testing.T) {
	ch := make(chan string, 30)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch <- r.Method
	}))
	collector := NewCollector(srv.URL, 2*time.Second, 10*time.Second, "")
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

	collector.stat = stat
	collector.sendReport()
	time.Sleep(2 * time.Second)
	requestCount := 0
	isAllRequestPost := true

	for len(ch) > 0 {
		method := <-ch
		if method != "POST" {
			isAllRequestPost = false
		}
		requestCount++
	}
	assert.Equal(t, 1, requestCount)
	assert.Equal(t, true, isAllRequestPost)
}
