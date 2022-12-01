package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type JSONValueHAndler struct {
	Server
}

func (jvh JSONValueHAndler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	val := jvh.storage.GetByKey(metr.ID)
	if val == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if val.GetShortTypeName() == metr.MType {
			respBody, err := json.Marshal(val.GetMetricEntity())
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.Write(respBody)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func NewJSONValueHAndler(s storage.Repository) *JSONValueHAndler {
	return &JSONValueHAndler{
		Server: Server{
			storage: s,
		},
	}
}
