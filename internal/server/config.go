package server

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/viking311/monitoring/internal/logger"
)

const (
	DefaultAddress       = "localhost:8080"
	DefaultStoreInterval = 300 * time.Second
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
	DefaultHashKey       = ""
	DefaultDatabaseDsn   = ""
)

type ServerConfig struct {
	Address       *string        `env:"ADDRESS"`
	StoreInterval *time.Duration `env:"STORE_INTERVAL"`
	StoreFile     *string        `env:"STORE_FILE"`
	Restore       *bool          `env:"RESTORE"`
	HashKey       *string        `env:"KEY"`
	DatabaseDsn   *string        `env:"DATABASE_DSN"`
}

var Config ServerConfig

func init() {
	logger.Logger.Info("start reading configuration")
	logger.Logger.Debug("reading flags")

	addressFlag := flag.String("a", DefaultAddress, "address to listen")
	restoreFlag := flag.Bool("r", DefaultRestore, "restore data from file")
	storeInterval := flag.Duration("i", DefaultStoreInterval, "how often store data to file")
	storeFile := flag.String("f", DefaultStoreFile, "name of file for storing")
	hashKey := flag.String("k", DefaultHashKey, "hash key")
	dbDsn := flag.String("d", DefaultDatabaseDsn, "connection to db")
	flag.Parse()

	logger.Logger.Debug("reading enviroments")
	if err := env.Parse(&Config); err != nil {
		logger.Logger.Fatal(err)
	}

	if Config.Address == nil {
		Config.Address = addressFlag
	}

	if Config.Restore == nil {
		Config.Restore = restoreFlag
	}

	if Config.StoreFile == nil {
		Config.StoreFile = storeFile
	}

	if Config.StoreInterval == nil {
		Config.StoreInterval = storeInterval
	}

	if Config.HashKey == nil {
		Config.HashKey = hashKey
	}

	if Config.DatabaseDsn == nil {
		Config.DatabaseDsn = dbDsn
	}

	logger.Logger.Info("finish reading configuration")
}
