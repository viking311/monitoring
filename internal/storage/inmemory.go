package storage

import (
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
}

func (ims *InMemoryStorage) Update(value entity.MetricEntityInterface) {
	item := ims.GetByKey(value.GetKey())
	if item != nil {
		item.SetValue(value.GetValue())
		ims.data[value.GetKey()] = item
	} else {
		ims.data[value.GetKey()] = value
	}
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

func (ims *InMemoryStorage) GetAll() []entity.MetricEntityInterface {
	slice := make([]entity.MetricEntityInterface, len(ims.data))
	for _, item := range ims.data {
		slice = append(slice, item)
	}

	return slice
}

func NewInMemoryStorage() InMemoryStorage {
	return InMemoryStorage{
		data: make(map[string]entity.MetricEntityInterface),
	}
}
