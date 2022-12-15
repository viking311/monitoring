package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/viking311/monitoring/internal/entity"
)

type DBStorage struct {
	db     *sql.DB
	upChan UpdateChannel
}

func (dbs *DBStorage) Update(value entity.Metrics) error {
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
		return err
	}

	select {
	case dbs.upChan <- struct{}{}:
	default:
	}

	return nil
}

func (dbs *DBStorage) Delete(key string) error {
	_, err := dbs.db.Exec("DELETE FROM metrics WHERE mkey=$1", key)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) GetByKey(key string) (entity.Metrics, error) {
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

func (dbs *DBStorage) BatchUpdate(values []entity.Metrics) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			log.Println()
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO metrics VALUES($1,$2,$3,$4,$5) ON CONFLICT (mkey) DO UPDATE SET delta=metrics.delta + $6, value=$7")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	for _, value := range values {
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

		if _, err := stmt.Exec(value.GetKey(), value.ID, value.MType, delta, floatValue, delta, floatValue); err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	select {
	case dbs.upChan <- struct{}{}:
	default:
	}

	return nil
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
		upChan: make(UpdateChannel),
	}, nil
}
