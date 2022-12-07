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
	defer r.Body.Close()

	var metr entity.Metrics

	err = json.Unmarshal(body, &metr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := jvh.storage.GetByKey(metr.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if val.MType == metr.MType {
			val.Hash = entity.MetricsHash(val, jvh.hashKey)
			respBody, err := json.Marshal(val)
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

func NewJSONValueHandler(s storage.Repository, hashKey string) *JSONValueHAndler {
	return &JSONValueHAndler{
		Server: Server{
			storage: s,
			hashKey: hashKey,
		},
	}
}
