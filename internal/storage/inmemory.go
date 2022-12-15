package storage

import (
	"fmt"
	"strings"
	"sync"

	"github.com/viking311/monitoring/internal/entity"
)

type InMemoryStorage struct {
	data         map[string]entity.Metrics
	mx           sync.RWMutex
	upChan       UpdateChannel
	isSendNotify bool
}

func (ims *InMemoryStorage) Update(value entity.Metrics) error {
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

	if ims.isSendNotify {
		go func() {
			ims.upChan <- struct{}{}
		}()
	}

	return nil
}

func (ims *InMemoryStorage) Delete(key string) error {
	ims.mx.Lock()
	defer ims.mx.Unlock()
	delete(ims.data, key)

	return nil
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

func (ims *InMemoryStorage) GetAll() ([]entity.Metrics, error) {
	ims.mx.RLock()
	defer ims.mx.RUnlock()

	slice := make([]entity.Metrics, 0, len(ims.data))
	for _, item := range ims.data {

		if item.Delta != nil {
			delta := *item.Delta
			item.Delta = &delta
		}

		if item.Value != nil {
			value := *item.Value
			item.Value = &value
		}

		slice = append(slice, item)
	}

	return slice, nil
}

func (ims *InMemoryStorage) GetUpdateChannal() UpdateChannel {
	return ims.upChan
}

func (ims *InMemoryStorage) BatchUpdate(values []entity.Metrics) error {
	ims.mx.Lock()
	defer ims.mx.Unlock()

	for _, value := range values {
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

	select {
	case ims.upChan <- struct{}{}:
	default:
	}

	return nil
}

func NewInMemoryStorage(isSendNotify bool) *InMemoryStorage {
	return &InMemoryStorage{
		data:         make(map[string]entity.Metrics),
		mx:           sync.RWMutex{},
		upChan:       make(UpdateChannel),
		isSendNotify: isSendNotify,
	}
}
