package main

import (
	"net/http"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/handlers"
	"github.com/viking311/monitoring/internal/storage"
)

func main() {
	c := make(chan entity.MetricEntityInterface, 100)

	saver := storage.NewSaver(c)
	go saver.Go()

	http.Handle("/update/", handlers.NewUpdateHandler(c))

	http.ListenAndServe(":8080", nil)
}
