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

type GetJsonValueHandler struct {
	storage storage.Repository
}

func (gjvh GetJsonValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	val := gjvh.storage.GetByKey(strings.ToLower(metricEntity.ID))

	if val == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if val.GetShortTypeName() == metricEntity.MType {
			w.Header().Add("Content-Type", "application/json")
			body, err := json.Marshal(val.GetMetricsEntity())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(body)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func NewGetJsonValueHandler(s storage.Repository) *GetJsonValueHandler {
	return &GetJsonValueHandler{
		storage: s,
	}
}
