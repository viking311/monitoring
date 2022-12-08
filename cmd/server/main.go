package main

//CREATE TABLE IF NOT EXISTS metrics (metric_key varchar(50) NOT NULL, metric_id varchar(50) NOT NULL, metric_type varchar(50) NOT NULL, metric_delta INT, metric_value DOUBLE PRECISION)
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
	var db *sql.DB

	s := storage.NewInMemoryStorage()

	if len(*server.Config.DatabaseDsn) > 0 {
		db_tmp, err := sql.Open("postgres", *server.Config.DatabaseDsn)
		db = db_tmp
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		sw, err := storage.NewSnapshotDbWriter(db, s, *server.Config.StoreInterval)
		if err != nil {
			log.Fatal(err)
		}

		defer sw.Close()

		if *server.Config.Restore {
			sw.Load()
		}

		go sw.Receive()
	}

	if len(*server.Config.StoreFile) > 0 && db == nil {
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
