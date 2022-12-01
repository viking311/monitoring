package storage

import (
	"strings"
	"sync"

	"github.com/viking311/monitoring/internal/entity"
)

type Repository interface {
	Update(value entity.MetricEntityInterface)
	Delete(key string)
	GetByKey(key string) entity.MetricEntityInterface
	GetAll() []entity.MetricEntityInterface
}

type InMemoryStorage struct {
	data map[string]entity.MetricEntityInterface
	mx   sync.RWMutex
}

func (ims *InMemoryStorage) Update(value entity.MetricEntityInterface) {
	ims.mx.Lock()
	defer ims.mx.Unlock()
	item, ok := ims.data[value.GetKey()]
	if ok {
		item.SetValue(value.GetValue())
		ims.data[value.GetKey()] = item
	} else {
		ims.data[value.GetKey()] = value
	}
}

func (ims *InMemoryStorage) Delete(key string) {
	ims.mx.Lock()
	defer ims.mx.Unlock()
	delete(ims.data, key)
}

func (ims *InMemoryStorage) GetByKey(key string) entity.MetricEntityInterface {
	ims.mx.RLock()
	defer ims.mx.RUnlock()
	value, ok := ims.data[strings.ToLower(key)]
	if ok {
		return value
	} else {
		return nil
	}
}

func (ims *InMemoryStorage) GetAll() []entity.MetricEntityInterface {
	ims.mx.RLock()
	defer ims.mx.RUnlock()
	slice := make([]entity.MetricEntityInterface, len(ims.data))
	for _, item := range ims.data {
		slice = append(slice, item)
	}

	return slice
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]entity.MetricEntityInterface),
		mx:   sync.RWMutex{},
	}
}
