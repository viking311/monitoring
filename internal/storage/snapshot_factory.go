package storage

import (
	"time"
)

func NewSpashotWriter(s Repository, fileName string, storeInterval time.Duration) (SnapshotWriterInterface, error) {
	if len(fileName) > 0 {
		sw, err := NewSnapshoFiletWriter(s, fileName, storeInterval)
		if err != nil {
			return nil, err
		}
		return sw, nil
	}
	return nil, nil
}
