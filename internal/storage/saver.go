package storage

import (
	"github.com/viking311/monitoring/internal/entity"
)

type Saver struct {
	storage    Repository
	updateChan chan entity.MetricEntityInterface
}

func (s Saver) Go() {
	for val := range s.updateChan {
		s.storage.Update(val)

		// fmt.Println("Updated metric - ", val.GetKey(), "with value - ", val.GetStringValue())
	}
}

func NewSaver(c chan entity.MetricEntityInterface, s Repository) Saver {

	return Saver{
		storage:    s,
		updateChan: c,
	}
}
