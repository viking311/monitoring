package storage

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viking311/monitoring/internal/entity"
)

func TestInMemoryStorage_Update_gauge(t *testing.T) {
	type args struct {
		value entity.Metrics
	}
	floatVal := float64(1)
	tests := []struct {
		name string
		ims  *InMemoryStorage
		args args
		want entity.Metrics
	}{
		{
			name: "InMemoryStorage_Update_test1",
			ims: &InMemoryStorage{
				data: make(map[string]entity.Metrics),
				mx:   sync.RWMutex{},
			},
			args: args{
				value: entity.Metrics{
					ID:    "m1",
					MType: "gauge",
					Value: &floatVal,
				},
			},
			want: entity.Metrics{
				ID:    "m1",
				MType: "gauge",
				Value: &floatVal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ims.Update(tt.args.value)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(tt.ims.data))
			value, ok := tt.ims.data[tt.args.value.GetKey()]
			assert.Equal(t, true, ok)
			assert.Equal(t, tt.want, value)
		})
	}
}

func TestInMemoryStorage_Update_counter(t *testing.T) {
	type args struct {
		value entity.Metrics
	}
	intVal := uint64(1)
	wantValue := uint64(2)
	tests := []struct {
		name string
		ims  *InMemoryStorage
		args []args
		want entity.Metrics
	}{
		{
			name: "InMemoryStorage_Update_test1",
			ims: &InMemoryStorage{
				data: make(map[string]entity.Metrics),
				mx:   sync.RWMutex{},
			},
			args: []args{
				{
					value: entity.Metrics{
						ID:    "m1",
						MType: "counter",
						Delta: &intVal,
					},
				},
				{
					value: entity.Metrics{
						ID:    "m1",
						MType: "counter",
						Delta: &intVal,
					},
				},
			},
			want: entity.Metrics{
				ID:    "m1",
				MType: "counter",
				Delta: &wantValue,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, arg := range tt.args {
				err := tt.ims.Update(arg.value)
				assert.Nil(t, err)
			}
			assert.Equal(t, 1, len(tt.ims.data))
			value, ok := tt.ims.data[tt.want.GetKey()]
			assert.Equal(t, true, ok)
			assert.Equal(t, tt.want, value)
		})
	}
}

func TestInMemoryStorage_Delete(t *testing.T) {
	type args struct {
		key string
	}

	ims := NewInMemoryStorage(false)

	ims.data["m1"] = entity.Metrics{
		ID:    "m1",
		MType: "gauge",
	}
	tests := []struct {
		name string
		ims  *InMemoryStorage
		args args
	}{
		{
			name: "TestInMemoryStorage_Delete",
			ims:  ims,
			args: args{
				key: "m1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ims.Delete(tt.args.key)
			assert.Empty(t, tt.ims.data)
		})
	}
}

func TestInMemoryStorage_GetByKey(t *testing.T) {
	type args struct {
		key string
	}
	type wants struct {
		found error
		value entity.Metrics
	}

	ims := NewInMemoryStorage(false)
	ims.data["m1"] = entity.Metrics{
		ID:    "m1",
		MType: "gauge",
	}
	tests := []struct {
		name string
		ims  *InMemoryStorage
		args args
		want wants
	}{
		{
			name: "TestInMemoryStorage_GetByKey_test1",
			ims:  ims,
			args: args{
				key: "m1",
			},
			want: wants{
				found: nil,
				value: ims.data["m1"],
			},
		},
		{
			name: "TestInMemoryStorage_GetByKey_test1",
			ims:  ims,
			args: args{
				key: "m2",
			},
			want: wants{
				found: fmt.Errorf("metric not found"),
				value: entity.Metrics{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.ims.GetByKey(tt.args.key)
			assert.Equal(t, tt.want.found, err)
			assert.Equal(t, tt.want.value, value)
		})
	}
}

func TestInMemoryStorage_getAll(t *testing.T) {
	intValue := uint64(1)
	floatValue := float64(1)
	m1 := entity.Metrics{
		ID:    "m1",
		MType: "counter",
		Delta: &intValue,
	}
	m2 := entity.Metrics{
		ID:    "m2",
		MType: "gauge",
		Value: &floatValue,
	}
	ims := NewInMemoryStorage(false)
	ims.data["m1"] = m1
	ims.data["m2"] = m2
	slice := make([]entity.Metrics, 0, len(ims.data))
	slice = append(slice, m1)
	slice = append(slice, m2)
	tests := []struct {
		name string
		ims  *InMemoryStorage
		want []entity.Metrics
	}{
		{
			name: "TestInMemoryStorage_getAll_test1",
			ims:  ims,
			want: slice,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := tt.ims.GetAll()
			assert.Nil(t, err)
			assert.Equal(t, tt.want, values)
		})
	}
}

func TestInMemoryStorage_butchUpdate(t *testing.T) {
	intValue := uint64(1)
	intValue2 := uint64(2)
	floatValue := float64(1)
	m1 := entity.Metrics{
		ID:    "m1",
		MType: "counter",
		Delta: &intValue,
	}
	m2 := entity.Metrics{
		ID:    "m2",
		MType: "gauge",
		Value: &floatValue,
	}
	m3 := entity.Metrics{
		ID:    "m1",
		MType: "counter",
		Delta: &intValue2,
	}

	ims := NewInMemoryStorage(false)

	slice := make([]entity.Metrics, 0, len(ims.data))
	slice = append(slice, m1)
	slice = append(slice, m2)
	slice = append(slice, m1)

	wantSlice := make([]entity.Metrics, 0, len(ims.data))
	wantSlice = append(wantSlice, m1)
	wantSlice = append(wantSlice, m3)

	tests := []struct {
		name string
		ims  *InMemoryStorage
		args []entity.Metrics
		want []entity.Metrics
	}{
		{
			name: "TestInMemoryStorage_getAll_test1",
			ims:  ims,
			args: slice,
			want: wantSlice,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ims.BatchUpdate(tt.args)
			assert.Nil(t, err)
			for _, w := range tt.want {
				value, err := tt.ims.GetByKey(w.GetKey())
				assert.Equal(t, value, w)
				assert.Nil(t, err)
			}
		})
	}

}
