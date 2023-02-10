package storage

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/viking311/monitoring/internal/entity"
)

func TestNewDBStorage_DBIsNil(t *testing.T) {
	wantError := fmt.Errorf("db instance is needed")
	dbs, err := NewDBStorage(nil, false)
	assert.Nil(t, dbs)
	assert.Equal(t, wantError, err)
}

func TestNewDBStorage_FailedCreateTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	wantError := fmt.Errorf("some error")
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS metrics").WillReturnError(wantError)
	dbs, err := NewDBStorage(db, false)
	assert.Nil(t, dbs)
	assert.Equal(t, wantError, err)
}

func TestNewDBStorage_success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS metrics").WillReturnResult(sqlmock.NewResult(1, 1))
	dbs, err := NewDBStorage(db, false)
	assert.Nil(t, err)
	assert.NotNil(t, dbs)
}
func TestNewDBStorage_update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS metrics").WillReturnResult(sqlmock.NewResult(1, 1))
	dbs, _ := NewDBStorage(db, true)
	intValue := uint64(1)
	floatValue := float64(1)
	tests := []struct {
		name    string
		storage *DBStorage
		arg     entity.Metrics
		want    error
	}{
		{
			name:    "TestNewDBStorage_update_counter_success",
			storage: dbs,
			arg: entity.Metrics{
				ID:    "counter_mentric",
				MType: "counter",
				Delta: &intValue,
			},
			want: nil,
		},
		{
			name:    "TestNewDBStorage_update_counter_error",
			storage: dbs,
			arg: entity.Metrics{
				ID:    "counter_mentric",
				MType: "counter",
				Delta: &intValue,
			},
			want: fmt.Errorf("some error"),
		},
		{
			name:    "TestNewDBStorage_update_gauge_success",
			storage: dbs,
			arg: entity.Metrics{
				ID:    "gauge_mentric",
				MType: "gauge",
				Value: &floatValue,
			},
			want: nil,
		},
		{
			name:    "TestNewDBStorage_update_counter_error",
			storage: dbs,
			arg: entity.Metrics{
				ID:    "gauge_mentric",
				MType: "gauge",
				Value: &floatValue,
			},
			want: fmt.Errorf("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == nil {
				mock.ExpectExec("INSERT INTO metrics").WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO metrics").WillReturnError(tt.want)
			}
			err = tt.storage.Update(tt.arg)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestNewDBStorage_delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS metrics").WillReturnResult(sqlmock.NewResult(1, 1))
	dbs, _ := NewDBStorage(db, true)
	tests := []struct {
		name    string
		storage *DBStorage
		want    error
	}{
		{
			name:    "TestNewDBStorage_delete_success",
			storage: dbs,
			want:    nil,
		},
		{
			name:    "TestNewDBStorage_delete_error",
			storage: dbs,
			want:    fmt.Errorf("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == nil {
				mock.ExpectExec("DELETE FROM metrics").WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("DELETE FROM metrics").WillReturnError(tt.want)
			}
			err = tt.storage.Delete("key")
			assert.Equal(t, tt.want, err)
		})
	}
}
