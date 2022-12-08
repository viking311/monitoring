package storage

import (
	"database/sql"
	"time"
)

func NewSpashotWriter(db *sql.DB, s Repository, fileName string, storeInterval time.Duration) (SnapshotWriterInterface, error) {
	if db != nil {
		sdw, err := NewSnapshotDbWriter(db, s, storeInterval)
		if err != nil {
			return nil, err
		}

		return sdw, nil
	}

	if len(fileName) > 0 {
		sw, err := NewSnapshoFiletWriter(s, fileName, storeInterval)
		if err != nil {
			return nil, err
		}
		return sw, nil
	}
	return nil, nil
}
