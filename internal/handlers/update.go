package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/viking311/monitoring/internal/entity"
)

type UpdateHandler struct {
	unpdateChan chan entity.MetricEntityInterface
}

func (uh UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch pathParts[1] {
	case "gauge":
		mValue, err := strconv.ParseFloat(pathParts[3], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		entity := entity.GaugeMetricEntity{
			Name:  pathParts[2],
			Value: mValue,
		}
		uh.unpdateChan <- &entity

	case "counter":
		mValue, err := strconv.ParseUint(pathParts[3], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		entity := entity.CounterMetricEntity{
			Name:  pathParts[2],
			Value: mValue,
		}
		uh.unpdateChan <- &entity
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

}

func NewUpdateHandler(c chan entity.MetricEntityInterface) UpdateHandler {
	return UpdateHandler{
		unpdateChan: c,
	}
}
