package storage

import (
	"database/sql"

	"github.com/viking311/monitoring/internal/server"
)

func NewSnapshotInstance(config *server.ServerConfig, db *sql.DB, store Repository) (SnapshotWriterInterface, error) {
	if len(*config.DatabaseDsn) > 0 {
		sdw, err := NewSnapshotDbWriter(db, store, *config.StoreInterval)
		if err != nil {
			return nil, err
		}
		return sdw, nil
	} else if len(*server.Config.StoreFile) > 0 {
		sw, err := NewSnapshotWriter(store, *server.Config.StoreFile, *config.StoreInterval)
		if err != nil {
			return nil, err
		}
		return sw, nil
	}

	return nil, nil
}
