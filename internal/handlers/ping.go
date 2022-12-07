package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type PingHandler struct {
	Server
}

func (ph PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ph.db == nil {
		fmt.Println("Db nuil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := ph.db.PingContext(ctx); err != nil {
		log.Println(err)
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
