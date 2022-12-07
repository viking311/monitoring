package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/viking311/monitoring/internal/handlers"
	"github.com/viking311/monitoring/internal/server"
	"github.com/viking311/monitoring/internal/storage"
)

func main() {

	db, err := sql.Open("postgres", *server.Config.DatabaseDsn)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

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
	r.Use(server.Gzip)
	r.Use(server.UnGzip)

	getListHandler := handlers.NewGetListHandler(s)
	r.Get("/", getListHandler.ServeHTTP)

	updateHandler := handlers.NewUpdatePlainTextHandler(s)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)

	jsonUpdateHandler := handlers.NewJSONUpdateHandler(s, *server.Config.HashKey)
	r.Post("/update/", jsonUpdateHandler.ServeHTTP)

	valueHandler := handlers.NewGetValueHandler(s)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)

	jsonValueHandler := handlers.NewJSONValueHandler(s, *server.Config.HashKey)
	r.Post("/value/", jsonValueHandler.ServeHTTP)

	pingHandler := handlers.NewPingHandler(db)
	r.Get("/ping", pingHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(*server.Config.Address, r))
}
