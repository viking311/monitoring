package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/viking311/monitoring/internal/handlers"
	"github.com/viking311/monitoring/internal/storage"
)

func main() {
	s := storage.NewInMemoryStorage()

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	getListHandler := handlers.NewGetListHandler(s)

	r.Get("/", getListHandler.ServeHTTP)

	updateHandler := handlers.NewUpdateHandler(s)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)

	jsonUpdateHandler := handlers.NewJSONUpdateHandler(s)
	r.Post("/update/", jsonUpdateHandler.ServeHTTP)

	valueHandler := handlers.NewGetValueHandler(s)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(":8080", r))
}
