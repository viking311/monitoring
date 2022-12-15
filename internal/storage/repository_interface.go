package storage

import "github.com/viking311/monitoring/internal/entity"

type UpdateChannel chan struct{}

type Repository interface {
	Update(value entity.Metrics) error
	Delete(key string) error
	GetByKey(key string) (entity.Metrics, error)
	GetAll() []entity.Metrics
	GetUpdateChannal() UpdateChannel
	BatchUpdate([]entity.Metrics) error
}
