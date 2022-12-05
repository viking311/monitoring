package server

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	DEFAULT_ADDRESS        = "localhost:8080"
	DEFAULT_STORE_INTERVAL = 300 * time.Second
	DEFAULT_STORE_FILE     = "/tmp/devops-metrics-db.json"
	DEFAULT_RESTORE        = true
)

type ServerConfig struct {
	Address       *string        `env:"ADDRESS"`
	StoreInterval *time.Duration `env:"STORE_INTERVAL"`
	StoreFile     *string        `env:"STORE_FILE"`
	Restore       *bool          `env:"RESTORE"`
}

var Config ServerConfig

func init() {
	addressFlag := flag.String("a", DEFAULT_ADDRESS, "address to listen")
	restoreFlag := flag.Bool("r", DEFAULT_RESTORE, "restore data from file")
	storeInterval := flag.Duration("i", DEFAULT_STORE_INTERVAL, "how often store data to file")
	storeFile := flag.String("f", DEFAULT_STORE_FILE, "name of file for storing")
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
