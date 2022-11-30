package handlers

import (
	"net/http"

	"github.com/viking311/monitoring/internal/storage"
)

type Server struct {
	storage storage.Repository
}

func (srv Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
