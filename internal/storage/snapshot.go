package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/viking311/monitoring/internal/entity"
	"github.com/viking311/monitoring/internal/logger"
)

type SnapshotWriter struct {
	store         Repository
	file          *os.File
	storeInterval time.Duration
	mx            sync.Mutex
	writer        bufio.Writer
}

func (sw *SnapshotWriter) Load() error {
	scanner := bufio.NewScanner(sw.file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		metric := entity.Metrics{}
		err := json.Unmarshal(scanner.Bytes(), &metric)
		if err != nil {
			return err
		}
		err = sw.store.Update(metric)
		if err != nil {
			return err
		}
	}
	logger.Info("data loaded from file " + sw.file.Name())
	return nil
}

func (sw *SnapshotWriter) Close() {
	err := sw.dump()
	if err != nil {
		logger.Error(err)
	}
	err = sw.file.Close()
	if err != nil {
		logger.Error(err)
	}
}

func (sw *SnapshotWriter) Receive() {
	if sw.storeInterval > 0 {
		ticker := time.NewTicker(sw.storeInterval)
		defer ticker.Stop()
		for range ticker.C {
			err := sw.dump()
			if err != nil {
				logger.Error(err)
			}
		}
	} else {
		for range sw.store.GetUpdateChannal() {
			err := sw.dump()
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (sw *SnapshotWriter) dump() error {
	sw.mx.Lock()
	defer sw.mx.Unlock()

	values, err := sw.store.GetAll()
	if err != nil {
		return err
	}
	for _, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = sw.writer.WriteString(string(data) + "\n")
		if err != nil {
			return err
		}
	}
	err = sw.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = sw.file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = sw.writer.Flush()
	if err != nil {
		return err
	}
	logger.Debug("data stored to file " + sw.file.Name())

	return nil
}

func NewSnapshotWriter(storage Repository, fileName string, storeInterval time.Duration) (*SnapshotWriter, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &SnapshotWriter{
		store:         storage,
		file:          file,
		storeInterval: storeInterval,
		mx:            sync.Mutex{},
		writer:        *bufio.NewWriter(file),
	}, nil
}
