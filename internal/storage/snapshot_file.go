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

type SnapshoFiletWriter struct {
	store         Repository
	file          *os.File
	storeInterval time.Duration
	mx            sync.Mutex
	writer        bufio.Writer
}

func (sw *SnapshoFiletWriter) Load() {
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

func (sw *SnapshoFiletWriter) Close() {
	sw.dump()
	sw.file.Close()
}

func (sw *SnapshoFiletWriter) Receive() {
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

func (sw *SnapshoFiletWriter) dump() {
	sw.mx.Lock()
	defer sw.mx.Unlock()

	for _, v := range sw.store.GetAll() {
		data, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
			continue
		}

		sw.writer.WriteString(string(data) + "\n")
	}
	sw.file.Truncate(0)
	sw.file.Seek(0, 0)
	sw.writer.Flush()
	log.Println("data stored to file " + sw.file.Name())
}

func NewSnapshoFiletWriter(storage Repository, fileName string, storeInterval time.Duration) (*SnapshoFiletWriter, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &SnapshoFiletWriter{
		store:         storage,
		file:          file,
		storeInterval: storeInterval,
		mx:            sync.Mutex{},
		writer:        *bufio.NewWriter(file),
	}, nil
}
