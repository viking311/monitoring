package handlers

import (
	"net/http"

	"github.com/viking311/monitoring/internal/storage"
)

type UpdateHandler struct {
	textPlainHandler http.Handler
	jsoHandler       http.Handler
}

func (uh UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "text/plain" {
		uh.textPlainHandler.ServeHTTP(w, r)
	} else if contentType == "application/json" {
		uh.jsoHandler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func NewUpdateHandler(s storage.Repository) *UpdateHandler {
	return &UpdateHandler{
		textPlainHandler: NewUpdatePlainTextHandler(s),
		jsoHandler:       NewJSONUpdateHandler(s),
	}
}
