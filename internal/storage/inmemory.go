package storage

import (
	"fmt"
	"strings"
	"sync"

	"github.com/viking311/monitoring/internal/entity"
)

type InMemoryStorage struct {
	data map[string]entity.Metrics
	mx   sync.RWMutex
}

func (ims *InMemoryStorage) Update(value entity.Metrics) {
	ims.mx.Lock()
	defer ims.mx.Unlock()

	if value.MType == "counter" {
		_, ok := ims.data[value.GetKey()]
		if ok && value.Delta != nil {
			*ims.data[value.GetKey()].Delta += *value.Delta
		} else if !ok && value.Delta != nil {
			ims.data[value.GetKey()] = value
		}
	} else if value.MType == "gauge" && value.Value != nil {
		ims.data[value.GetKey()] = value
	}
}

func (ims *InMemoryStorage) Delete(key string) {
	ims.mx.Lock()
	defer ims.mx.Unlock()
	delete(ims.data, key)
}

func (ims *InMemoryStorage) GetByKey(key string) (entity.Metrics, error) {
	ims.mx.RLock()
	defer ims.mx.RUnlock()

	value, ok := ims.data[strings.ToLower(key)]
	if ok {
		return value, nil
	} else {
		return entity.Metrics{}, fmt.Errorf("metric not found")
	}
}

func (ims *InMemoryStorage) GetAll() []entity.Metrics {
	ims.mx.RLock()
	defer ims.mx.RUnlock()

	slice := make([]entity.Metrics, len(ims.data))
	i := 0
	for _, item := range ims.data {
		slice[i] = item
		i++
	}

	return slice
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]entity.Metrics),
		mx:   sync.RWMutex{},
	}
}
