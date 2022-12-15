package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/logger"
	"github.com/viking311/monitoring/internal/storage"
)

type JSONBatchUpdateHandler struct {
	Server
}

func (jbuh *JSONBatchUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		logger.Logger.Error("incorrect content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var metricsCollection []entity.Metrics

	err = json.Unmarshal(body, &metricsCollection)
	if err != nil {
		logger.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cleanMetrics := make([]entity.Metrics, len(metricsCollection))
	i := 0
	for _, m := range metricsCollection {
		if jbuh.verifyMetricsSing(m) {
			cleanMetrics[i] = m
			i++
		}
	}

	err = jbuh.storage.BatchUpdate(cleanMetrics[:i])
	if err != nil {
		logger.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func NewJSONBatchUpdateHandler(s storage.Repository, hashKey string) *JSONBatchUpdateHandler {
	return &JSONBatchUpdateHandler{
		Server{
			storage: s,
			hashKey: hashKey,
		},
	}
}
