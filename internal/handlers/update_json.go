package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type JsonUpdateHandler struct {
	str storage.Repository
}

func (juh JsonUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentTypeHeader := r.Header.Get("Content-Type")

	if contentTypeHeader != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metricEntity entity.Metrics

	err = json.Unmarshal(body, &metricEntity)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metricEntity.MType {
	case "gauge":
		entity := entity.GaugeMetricEntity{
			Name:  strings.ToLower(metricEntity.ID),
			Value: *metricEntity.Value,
		}
		juh.str.Update(&entity)

	case "counter":
		entity := entity.CounterMetricEntity{
			Name:  strings.ToLower(metricEntity.ID),
			Value: *metricEntity.Delta,
		}
		juh.str.Update(&entity)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

}

func NewJsonUpdateHandler(s storage.Repository) *JsonUpdateHandler {
	return &JsonUpdateHandler{
		str: s,
	}
}
