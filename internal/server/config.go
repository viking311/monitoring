package server

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultAddress       = "localhost:8080"
	DefaultStoreInterval = 300 * time.Second
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
)

type ServerConfig struct {
	Address       *string        `env:"ADDRESS"`
	StoreInterval *time.Duration `env:"STORE_INTERVAL"`
	StoreFile     *string        `env:"STORE_FILE"`
	Restore       *bool          `env:"RESTORE"`
}

var Config ServerConfig

func init() {
	addressFlag := flag.String("a", DefaultAddress, "address to listen")
	restoreFlag := flag.Bool("r", DefaultRestore, "restore data from file")
	storeInterval := flag.Duration("i", DefaultStoreInterval, "how often store data to file")
	storeFile := flag.String("f", DefaultStoreFile, "name of file for storing")
	flag.Parse()

	if err := env.Parse(&Config); err != nil {
		log.Fatal(err)
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
}
