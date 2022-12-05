package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/viking311/monitoring/internal/handlers"
	"github.com/viking311/monitoring/internal/server"
	"github.com/viking311/monitoring/internal/storage"
)

func main() {

	s := storage.NewInMemoryStorage()

	if len(*server.Config.StoreFile) > 0 {
		sw, err := storage.NewSnapshotWriter(s, *server.Config.StoreFile, *server.Config.StoreInterval)
		if err != nil {
			log.Fatal(err)
		}
		defer sw.Close()

		if *server.Config.Restore {
			sw.Load()
		}

		go sw.Receive()

	}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	getListHandler := handlers.NewGetListHandler(s)

	r.Get("/", getListHandler.ServeHTTP)

	updateHandler := handlers.NewUpdatePlainTextHandler(s)
	// updateHandler := handlers.NewUpdateHandler(s)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)

	jsonUpdateHandler := handlers.NewJSONUpdateHandler(s)
	r.Post("/update/", jsonUpdateHandler.ServeHTTP)

	valueHandler := handlers.NewGetValueHandler(s)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)

	jsonValueHandler := handlers.NewJSONValueHAndler(s)
	r.Post("/value/", jsonValueHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(*server.Config.Address, r))
}
