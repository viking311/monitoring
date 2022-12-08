package storage

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type SnapshotDbWriter struct {
	db            *sql.DB
	store         Repository
	storeInterval time.Duration
	mx            sync.Mutex
}

func (sdw *SnapshotDbWriter) Load() {
	log.Println("data loaded from db")
}

func (sdw *SnapshotDbWriter) Receive() {

}

func (sdw *SnapshotDbWriter) Close() {

}

func NewSnapshotDbWriter(db *sql.DB, store Repository, storeInterval time.Duration) (*SnapshotDbWriter, error) {
	return &SnapshotDbWriter{
		db:            db,
		store:         store,
		storeInterval: storeInterval,
		mx:            sync.Mutex{},
	}, nil
}
