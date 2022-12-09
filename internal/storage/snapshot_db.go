package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/viking311/monitoring/internal/entity"
)

type SnapshotDBWriter struct {
	db    *sql.DB
	store Repository
	mx    sync.Mutex
}

func (sdw *SnapshotDBWriter) Load() {
	rows, err := sdw.db.Query("SELECT id, mtype, delta, value FROM metrics")
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
	err = rows.Err()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("data loaded from db")
	}

}

func (sdw *SnapshotDBWriter) Receive() {
	for range sdw.store.GetUpdateChannal() {
		sdw.dump()
	}

}

func (sdw *SnapshotDBWriter) Close() {
	sdw.dump()
}

func (sdw *SnapshotDBWriter) dump() {
	sdw.mx.Lock()
	defer sdw.mx.Unlock()

	for _, v := range sdw.store.GetAll() {
		delta := sql.NullInt64{}
		if v.Delta != nil {
			delta.Int64 = int64(*v.Delta)
			delta.Valid = true

		}

		value := sql.NullFloat64{}
		if v.Value != nil {
			value.Float64 = *v.Value
			value.Valid = true
		}
		_, err := sdw.db.Exec("INSERT INTO metrics VALUES($1,$2,$3,$4,$5) ON CONFLICT (mkey,mtype) DO UPDATE SET delta=$6, value=$7", v.GetKey(), v.ID, v.MType, delta, value, delta, value)
		if err != nil {
			log.Println(err)
		}
	}

}
func NewSnapshotDBWriter(db *sql.DB, store Repository) (*SnapshotDBWriter, error) {
	if db == nil {
		return nil, fmt.Errorf("db instance is needed")
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS metrics (mkey TEXT NOT NULL, id TEXT NOT NULL, mtype TEXT NOT NULL, delta BIGINT, value DOUBLE PRECISION, CONSTRAINT key_type_uniq UNIQUE (mkey, mtype))")
	if err != nil {
		return nil, err
	}

	return &SnapshotDBWriter{
		db:    db,
		store: store,
		mx:    sync.Mutex{},
	}, nil
}
