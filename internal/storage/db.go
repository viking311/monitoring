package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/viking311/monitoring/internal/entity"
)

type DBStorage struct {
	db     *sql.DB
	mx     sync.RWMutex
	upChan UpdateChannel
}

func (dbs *DBStorage) Update(value entity.Metrics) {
	dbs.mx.Lock()
	defer dbs.mx.Unlock()
	delta := sql.NullInt64{}
	if value.Delta != nil {
		delta.Int64 = int64(*value.Delta)
		delta.Valid = true

	}

	floatValue := sql.NullFloat64{}
	if value.Value != nil {
		floatValue.Float64 = *value.Value
		floatValue.Valid = true
	}
	_, err := dbs.db.Exec("INSERT INTO metrics VALUES($1,$2,$3,$4,$5) ON CONFLICT (mkey) DO UPDATE SET delta=metrics.delta + $6, value=$7", value.GetKey(), value.ID, value.MType, delta, floatValue, delta, floatValue)
	if err != nil {
		log.Println(err)
	}

	select {
	case dbs.upChan <- struct{}{}:
	default:
	}
}

func (dbs *DBStorage) Delete(key string) {
	dbs.mx.Lock()
	defer dbs.mx.Unlock()
	_, err := dbs.db.Exec("DELETE FROM metrics WHERE mkey=$1", key)
	if err != nil {
		log.Println(err)
	}
}

func (dbs *DBStorage) GetByKey(key string) (entity.Metrics, error) {
	dbs.mx.RLock()
	defer dbs.mx.RUnlock()

	metric := entity.Metrics{}

	rows, err := dbs.db.Query("SELECT id, mtype, delta, value FROM metrics WHERE mkey=$1", key)
	if err != nil {
		log.Println(err)
		return metric, err
	}
	for rows.Next() {
		var (
			delta sql.NullInt64
			value sql.NullFloat64
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

		err = rows.Err()
		if err != nil {
			log.Println(err)
			return metric, err
		}
	}

	return metric, nil
}

func (dbs *DBStorage) GetAll() []entity.Metrics {
	dbs.mx.RLock()
	defer dbs.mx.RUnlock()

	var count uint64
	err := dbs.db.QueryRow("SELECT COUNT(*) FROM metrics").Scan(&count)
	if err != nil {
		log.Println(err)
		return []entity.Metrics{}
	}

	if count == 0 {
		return []entity.Metrics{}
	}

	slice := make([]entity.Metrics, count)
	i := 0

	rows, err := dbs.db.Query("SELECT id, mtype, delta, value FROM metrics")
	if err != nil {
		log.Println(err)
		return []entity.Metrics{}
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
		slice[i] = metric
		i++
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return slice
}

func (dbs *DBStorage) GetUpdateChannal() UpdateChannel {
	return dbs.upChan
}

func NewDBStorage(db *sql.DB) (*DBStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("db instance is needed")
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS metrics (mkey TEXT NOT NULL PRIMARY KEY, id TEXT NOT NULL, mtype TEXT NOT NULL, delta BIGINT, value DOUBLE PRECISION)")
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		db:     db,
		mx:     sync.RWMutex{},
		upChan: make(UpdateChannel),
	}, nil
}
