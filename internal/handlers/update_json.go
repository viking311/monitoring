package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type JSONUpdateHandler struct {
	Server
}

func (juh JSONUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metr entity.Metrics

	err = json.Unmarshal(body, &metr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metr.MType {
	case "gauge":
		entity := entity.GaugeMetricEntity{
			Name:  metr.ID,
			Value: *metr.Value,
		}
		juh.storage.Update(&entity)

	case "counter":
		entity := entity.CounterMetricEntity{
			Name:  metr.ID,
			Value: *metr.Delta,
		}
		juh.storage.Update(&entity)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	currentValue := juh.storage.GetByKey(metr.ID)
	respBody, err := json.Marshal(currentValue.GetMetricEntity())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(respBody)
}

func NewJSONUpdateHandler(s storage.Repository) *JSONUpdateHandler {
	return &JSONUpdateHandler{
		Server: Server{
			storage: s,
		},
	}
}
