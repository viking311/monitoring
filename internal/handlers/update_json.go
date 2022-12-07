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
	defer r.Body.Close()

	var metr entity.Metrics

	err = json.Unmarshal(body, &metr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !juh.verifyMetricsSing(metr) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metr.Hash = ""

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

	currentValue.Hash = entity.MetricsHash(currentValue, juh.hashKey)

	respBody, err := json.Marshal(currentValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(respBody)
}

func NewJSONUpdateHandler(s storage.Repository, hashKey string) *JSONUpdateHandler {
	return &JSONUpdateHandler{
		Server: Server{
			storage: s,
			hashKey: hashKey,
		},
	}
}
