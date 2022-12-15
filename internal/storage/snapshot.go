package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/viking311/monitoring/internal/entity"
)

type SnapshotWriter struct {
	store         Repository
	file          *os.File
	storeInterval time.Duration
	mx            sync.Mutex
	writer        bufio.Writer
}

func (sw *SnapshotWriter) Load() {
	scanner := bufio.NewScanner(sw.file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		metric := entity.Metrics{}
		err := json.Unmarshal(scanner.Bytes(), &metric)
		if err != nil {
			log.Println(err)
			continue
		}
		sw.store.Update(metric)
	}
	log.Println("data loaded from file " + sw.file.Name())
}

func (sw *SnapshotWriter) Close() {
	sw.dump()
	sw.file.Close()
}

func (sw *SnapshotWriter) Receive() {
	if sw.storeInterval > 0 {
		ticker := time.NewTicker(sw.storeInterval)
		defer ticker.Stop()
		for range ticker.C {
			sw.dump()
		}
	} else {
		for range sw.store.GetUpdateChannal() {
			sw.dump()
		}
	}
}

func (sw *SnapshotWriter) dump() {
	sw.mx.Lock()
	defer sw.mx.Unlock()
	values, err := sw.store.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = sw.writer.WriteString(string(data) + "\n")
		if err != nil {
			log.Println(err)
		}
	}
	err = sw.file.Truncate(0)
	if err != nil {
		log.Println(err)
	}
	_, err = sw.file.Seek(0, 0)
	if err != nil {
		log.Println(err)
	}
	err = sw.writer.Flush()
	if err != nil {
		log.Println(err)
	}
	log.Println("data stored to file " + sw.file.Name())
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
