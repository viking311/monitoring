package handlers

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type GetValueHandler struct {
	storage storage.Repository
}

func (gvh GetValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	if typeName != "gauge" && typeName != "counter" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	valueName := strings.ToLower(chi.URLParam(r, "name"))

	val := gvh.storage.GetByKey(valueName)

	if val == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if checkType(typeName, val) {
			w.Header().Add("application-type", "text/plain")
			w.Write([]byte(val.GetStringValue()))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func NewGetValueHandler(s storage.Repository) GetValueHandler {
	return GetValueHandler{
		storage: s,
	}
}

func checkType(needle string, obj entity.MetricEntityInterface) bool {

	fullTypeName := strings.TrimPrefix(reflect.TypeOf(obj).String(), "*")
	shortName := ""
	if fullTypeName == "entity.CounterMetricEntity" {
		shortName = "counter"
	} else if fullTypeName == "entity.GaugeMetricEntity" {
		shortName = "gauge"
	}

	return shortName == needle
}
