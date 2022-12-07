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

func (juh *JSONUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if metr.MType != "gauge" && metr.MType != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	juh.storage.Update(metr)

	currentValue, err := juh.storage.GetByKey(metr.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	respBody, err := json.Marshal(currentValue)
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
