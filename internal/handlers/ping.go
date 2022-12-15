package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/viking311/monitoring/internal/logger"
)

type PingHandler struct {
	Server
}

func (ph *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ph.db == nil {
		logger.Logger.Debug("db is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := ph.db.PingContext(ctx); err != nil {
		logger.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NewPingHandler(db *sql.DB) *PingHandler {
	return &PingHandler{
		Server: Server{
			db: db,
		},
	}
}
