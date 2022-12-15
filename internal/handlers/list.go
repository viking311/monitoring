package handlers

import (
	"log"
	"net/http"

	"github.com/viking311/monitoring/internal/storage"
)

type GetListHandler struct {
	Server
}

func (glh GetListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := "<html><body><table border=1><tr><th>Metric</th><th>Value</th></tr>"
	values, err := glh.storage.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, v := range values {
		body += "<tr><td>" + v.ID + "</td><td>" + v.GetStringValue() + "</td></tr>"
	}
	body += "</table></body></html>"
	w.Header().Add("Content-Type", "text/html")
	_, err = w.Write([]byte(body))
	if err != nil {
		log.Println(err)
	}
}

func NewGetListHandler(s storage.Repository) *GetListHandler {
	return &GetListHandler{
		Server: Server{
			storage: s,
		},
	}
}
