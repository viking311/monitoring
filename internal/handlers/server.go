package handlers

import (
	"crypto/hmac"
	"database/sql"
	"net/http"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/storage"
)

type Server struct {
	storage storage.Repository
	hashKey string
	db      *sql.DB
}

func (srv Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (srv Server) verifyMetricsSing(data entity.Metrics) bool {
	if len(srv.hashKey) > 0 {
		hash := entity.MetricsHash(data, srv.hashKey)
		if !hmac.Equal([]byte(hash), []byte(data.Hash)) {
			return false
		}
	}
	return true
}
