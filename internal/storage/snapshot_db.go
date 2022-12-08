package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/viking311/monitoring/internal/entity"
)

type SnapshotDbWriter struct {
	db            *sql.DB
	store         Repository
	storeInterval time.Duration
	mx            sync.Mutex
}

func (sdw *SnapshotDbWriter) Load() {
	rows, err := sdw.db.Query("SELECT metric_id, metric_type,metric_delta, metric_value FROM metrics")
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var (
			metric entity.Metrics
			delta  sql.NullInt64
			value  sql.NullFloat64
		)

		err := rows.Scan(&metric.ID, &metric.MType, &delta, &value)
		if err != nil {
			log.Println(err)
			continue
		}
		if delta.Valid {
			val64 := uint64(delta.Int64)
			metric.Delta = &val64
		}

		if value.Valid {
			metric.Value = &value.Float64
		}
		sdw.store.Update(metric)
	}
	log.Println("data loaded from db")
}

func (sdw *SnapshotDbWriter) Receive() {
	if sdw.storeInterval > 0 {
		ticker := time.NewTicker(sdw.storeInterval)
		defer ticker.Stop()
		for range ticker.C {
			sdw.dump()
		}
	} else {
		for range sdw.store.GetUpdateChannal() {
			sdw.dump()
		}
	}
}

func (sdw *SnapshotDbWriter) dump() {
	sdw.mx.Lock()
	defer sdw.mx.Unlock()
	_, err := sdw.db.Exec("DELETE FROM metrics")

	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range sdw.store.GetAll() {
		var delta sql.NullInt64 = sql.NullInt64{}
		if v.Delta != nil {
			delta.Int64 = int64(*v.Delta)
			delta.Valid = true

		}

		var value sql.NullFloat64 = sql.NullFloat64{}
		if v.Value != nil {
			value.Float64 = *v.Value
			value.Valid = true
		}
		_, err := sdw.db.Exec("INSERT INTO metrics VALUES($1,$2,$3,$4,$5)", v.GetKey(), v.ID, v.MType, delta, value)
		if err != nil {
			log.Println(err)
		}
	}

}

func (sdw *SnapshotDbWriter) Close() {
	sdw.dump()
}

func NewSnapshotDbWriter(db *sql.DB, store Repository, storeInterval time.Duration) (*SnapshotDbWriter, error) {
	if db == nil {
		return nil, fmt.Errorf("db instance is needed")
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS metrics (metric_key varchar(50) NOT NULL, metric_id varchar(50) NOT NULL, metric_type varchar(50) NOT NULL, metric_delta INT, metric_value DOUBLE PRECISION)")
	if err != nil {
		return nil, err
	}

	return &SnapshotDbWriter{
		db:            db,
		store:         store,
		mx:            sync.Mutex{},
		storeInterval: storeInterval,
	}, nil
}
