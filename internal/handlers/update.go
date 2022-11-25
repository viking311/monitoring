package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type UpdateHandler struct {
	str storage.Repository
}

func (uh UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	typeName := strings.ToLower(chi.URLParam(r, "type"))
	metricName := strings.ToLower(chi.URLParam(r, "name"))
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
		uh.str.Update(&entity)

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
		uh.str.Update(&entity)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

}

func NewUpdateHandler(s storage.Repository) UpdateHandler {
	return UpdateHandler{
		str: s,
	}
}
