package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type UpdatePlainTextHandler struct {
	Server
}

func (uh UpdatePlainTextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")

	if contentType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	typeName := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if typeName == "" || metricName == "" || metricValue == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch typeName {
	case "gauge":
		mValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		entity := entity.GaugeMetricEntity{
			Name:  metricName,
			Value: mValue,
		}
		uh.storage.Update(&entity)

	case "counter":
		mValue, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		entity := entity.CounterMetricEntity{
			Name:  metricName,
			Value: mValue,
		}
		uh.storage.Update(&entity)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

}

func NewUpdatePlainTextHandler(s storage.Repository) *UpdatePlainTextHandler {
	return &UpdatePlainTextHandler{
		Server: Server{
			storage: s,
		},
	}
}
