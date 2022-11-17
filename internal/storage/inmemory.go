package storage

import (
	"github.com/viking311/monitoring/internal/entity"
)

type Repository interface {
	Update(value entity.MetricEntityInterface)
	Delete(key string)
	GetByKey(key string) entity.MetricEntityInterface
	getAll() []entity.MetricEntityInterface
}

type InMemoryStorage struct {
	data map[string]entity.MetricEntityInterface
}

func (ims *InMemoryStorage) Update(value entity.MetricEntityInterface) {
	ims.data[value.GetKey()] = value
}

func (ims *InMemoryStorage) Delete(key string) {
	delete(ims.data, key)
}

func (ims *InMemoryStorage) GetByKey(key string) entity.MetricEntityInterface {
	value, ok := ims.data[key]
	if ok {
		return value
	} else {
		return nil
	}
}

func (ims *InMemoryStorage) getAll() []entity.MetricEntityInterface {
	slice := make([]entity.MetricEntityInterface, len(ims.data))

	for _, item := range ims.data {
		slice = append(slice, item)
	}

	return slice
}

func MewInMemoryStorage() InMemoryStorage {
	return InMemoryStorage{
		data: make(map[string]entity.MetricEntityInterface),
	}
}
