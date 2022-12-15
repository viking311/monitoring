package handlers

import (
	"log"
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

func (uh *UpdatePlainTextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	typeName := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if typeName == "" || metricName == "" || metricValue == "" {
		log.Println("require params are empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if typeName != "gauge" && typeName != "counter" {
		log.Println("unknown metric type")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	metric := entity.Metrics{
		ID:    metricName,
		MType: typeName,
	}
	switch typeName {
	case "gauge":
		mValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Value = &mValue

	case "counter":
		mValue, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Delta = &mValue
	}
	err := uh.storage.Update(metric)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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
