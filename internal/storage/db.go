package storage

import (
	"database/sql"
	"fmt"

	"github.com/viking311/monitoring/internal/entity"
)

type DBStorage struct {
	db           *sql.DB
	upChan       UpdateChannel
	isSendNotify bool
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

	dbs.notify()

	return nil
}

func (dbs *DBStorage) Delete(key string) error {
	_, err := dbs.db.Exec("DELETE FROM metrics WHERE mkey=$1", key)
	if err != nil {
		return err
	}

	dbs.notify()

	return nil
}

func (dbs *DBStorage) GetByKey(key string) (entity.Metrics, error) {
	var (
		metric entity.Metrics
		delta  sql.NullInt64
		value  sql.NullFloat64
	)
	if err := dbs.db.QueryRow("SELECT id, mtype, delta, value FROM metrics WHERE mkey=$1", key).Scan(&metric.ID, &metric.MType, &delta, &value); err != nil {
		if err == sql.ErrNoRows {
			return metric, fmt.Errorf("not found metric with key '%s'", key)
		}
		return metric, err
	}
	if delta.Valid {
		val64 := uint64(delta.Int64)
		metric.Delta = &val64
	}

	if value.Valid {
		metric.Value = &value.Float64
	}
	return metric, nil
}

func (dbs *DBStorage) GetAll() ([]entity.Metrics, error) {
	var count uint64
	err := dbs.db.QueryRow("SELECT COUNT(*) FROM metrics").Scan(&count)
	if err != nil {
		return []entity.Metrics{}, err
	}

	if count == 0 {
		return []entity.Metrics{}, nil
	}

	slice := make([]entity.Metrics, 0, count)

	rows, err := dbs.db.Query("SELECT id, mtype, delta, value FROM metrics")
	if err != nil {
		return []entity.Metrics{}, err
	}

	for rows.Next() {
		var (
			metric entity.Metrics
			delta  sql.NullInt64
			value  sql.NullFloat64
		)

		err := rows.Scan(&metric.ID, &metric.MType, &delta, &value)
		if err != nil {
			return []entity.Metrics{}, err
		}
		if delta.Valid {
			val64 := uint64(delta.Int64)
			metric.Delta = &val64
		}

		if value.Valid {
			metric.Value = &value.Float64
		}
		slice = append(slice, metric)
	}
	err = rows.Err()
	if err != nil {
		return []entity.Metrics{}, err
	}

	return slice, nil
}

func (dbs *DBStorage) GetUpdateChannal() UpdateChannel {
	return dbs.upChan
}

func (dbs *DBStorage) BatchUpdate(values []entity.Metrics) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO metrics VALUES($1,$2,$3,$4,$5) ON CONFLICT (mkey) DO UPDATE SET delta=metrics.delta + $6, value=$7")
	if err != nil {
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
		return err
	}

	dbs.notify()

	return nil
}

func (dbs *DBStorage) notify() {
	if dbs.isSendNotify {
		go func() {
			dbs.upChan <- struct{}{}
		}()
	}
}

func NewDBStorage(db *sql.DB, isSendNotify bool) (*DBStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("db instance is needed")
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS metrics (mkey TEXT NOT NULL PRIMARY KEY, id TEXT NOT NULL, mtype TEXT NOT NULL, delta BIGINT, value DOUBLE PRECISION)")
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		db:           db,
		upChan:       make(UpdateChannel),
		isSendNotify: isSendNotify,
	}, nil
}
