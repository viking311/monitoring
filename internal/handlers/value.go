package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/logger"
	"github.com/viking311/monitoring/internal/storage"
)

type GetValueHandler struct {
	Server
}

func (gvh GetValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	if typeName != "gauge" && typeName != "counter" {
		logger.Logger.Error("unknown metric type")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	valueName := strings.ToLower(chi.URLParam(r, "name"))

	metric := entity.Metrics{
		ID:    valueName,
		MType: typeName,
	}
	val, err := gvh.storage.GetByKey(metric.GetKey())

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if val.MType == typeName {
			w.Header().Add("application-type", "text/plain")
			_, err := w.Write([]byte(val.GetStringValue()))
			if err != nil {
				logger.Logger.Error(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func NewGetValueHandler(s storage.Repository) *GetValueHandler {
	return &GetValueHandler{
		Server: Server{
			storage: s,
		},
	}
}
