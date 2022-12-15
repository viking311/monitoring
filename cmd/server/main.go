package main

import (
	"database/sql"
	"net/http"

	middlewareLogger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/viking311/monitoring/internal/handlers"
	"github.com/viking311/monitoring/internal/logger"
	"github.com/viking311/monitoring/internal/server"
	"github.com/viking311/monitoring/internal/storage"
)

var (
	db    *sql.DB
	sw    *storage.SnapshotWriter
	store storage.Repository
)

func main() {

	logger.Logger.Info("server start")

	defer logger.Logger.Info("server stop")

	if len(*server.Config.DatabaseDsn) > 0 {
		err := initDB(*server.Config.DatabaseDsn)
		if err != nil {
			logger.Logger.Fatal(err)
		}
	}
	isSendNotify := *server.Config.StoreInterval == 0
	if db == nil {
		store = storage.NewInMemoryStorage(isSendNotify)
	} else {
		var err error
		store, err = storage.NewDBStorage(db, isSendNotify)
		if err != nil {
			logger.Logger.Fatal(err)
		}
	}

	if (len(*server.Config.StoreFile) > 0) && db == nil {
		var err error

		sw, err = storage.NewSnapshotWriter(store, *server.Config.StoreFile, *server.Config.StoreInterval)
		if err != nil {
			logger.Logger.Fatal(err)
		}

		if *server.Config.Restore {
			err = sw.Load()
			if err != nil {
				logger.Logger.Error(err)
			}
		}

		go sw.Receive()

	}

	defer func() {
		if sw != nil {
			sw.Close()
		}
		if db != nil {
			db.Close()
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middlewareLogger.Logger("router", logger.Logger))
	r.Use(middleware.Recoverer)
	r.Use(server.Gzip)
	r.Use(server.UnGzip)

	getListHandler := handlers.NewGetListHandler(store)
	r.Get("/", getListHandler.ServeHTTP)

	updateHandler := handlers.NewUpdatePlainTextHandler(store)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)

	jsonUpdateHandler := handlers.NewJSONUpdateHandler(store, *server.Config.HashKey)
	r.Post("/update/", jsonUpdateHandler.ServeHTTP)

	valueHandler := handlers.NewGetValueHandler(store)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)

	jsonValueHandler := handlers.NewJSONValueHandler(store, *server.Config.HashKey)
	r.Post("/value/", jsonValueHandler.ServeHTTP)

	pingHandler := handlers.NewPingHandler(db)
	r.Get("/ping", pingHandler.ServeHTTP)

	jsonBatchUpdateHandler := handlers.NewJSONBatchUpdateHandler(store, *server.Config.HashKey)
	r.Post("/updates/", jsonBatchUpdateHandler.ServeHTTP)
	logger.Logger.Fatal(http.ListenAndServe(*server.Config.Address, r))
}

func initDB(dsn string) error {
	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	return db.Ping()
}
