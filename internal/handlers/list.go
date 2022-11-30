package handlers

import (
	"net/http"

	"github.com/viking311/monitoring/internal/storage"
)

type GetListHandler struct {
	Server
}

func (glh GetListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := "<html><body><table border=1><tr><th>Metric</th><th>Value</th></tr>"
	for _, v := range glh.storage.GetAll() {
		if v != nil {
			body += "<tr><td>" + v.GetKey() + "</td><td>" + v.GetStringValue() + "</td></tr>"
		}
	}
	body += "</table></body></html>"
	w.Header().Add("application-type", "text/html")
	w.Write([]byte(body))
}

func NewGetListHandler(s storage.Repository) *GetListHandler {
	return &GetListHandler{
		Server: Server{
			storage: s,
		},
	}
}
