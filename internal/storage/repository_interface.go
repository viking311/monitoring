package storage

import "github.com/viking311/monitoring/internal/entity"

type Repository interface {
	Update(value entity.Metrics)
	Delete(key string)
	GetByKey(key string) (entity.Metrics, error)
	GetAll() []entity.Metrics
}
