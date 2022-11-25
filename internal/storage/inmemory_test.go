package storage

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viking311/monitoring/internal/entity"
)

func TestInMemoryStorage_Update(t *testing.T) {
	type args struct {
		value entity.MetricEntityInterface
	}
	tests := []struct {
		name string
		ims  *InMemoryStorage
		args args
		want entity.MetricEntityInterface
	}{
		{
			name: "InMemoryStorage_Update_test1",
			ims: &InMemoryStorage{
				data: make(map[string]entity.MetricEntityInterface),
				mx:   sync.RWMutex{},
			},
			args: args{
				value: &entity.GaugeMetricEntity{
					Name:  "m1",
					Value: float64(1),
				},
			},
			want: &entity.GaugeMetricEntity{
				Name:  "m1",
				Value: float64(1),
			},
		},
		{
			name: "InMemoryStorage_Update_test2",
			ims: &InMemoryStorage{
				data: make(map[string]entity.MetricEntityInterface),
				mx:   sync.RWMutex{},
			},
			args: args{
				value: &entity.CounterMetricEntity{
					Name:  "m1",
					Value: uint64(1),
				},
			},
			want: &entity.CounterMetricEntity{
				Name:  "m1",
				Value: uint64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ims.Update(tt.args.value)
			assert.Equal(t, 1, len(tt.ims.data))
			value, ok := tt.ims.data[tt.args.value.GetKey()]
			assert.Equal(t, true, ok)
			assert.Equal(t, tt.want, value)
		})
	}
}

func TestInMemoryStorage_Delete(t *testing.T) {
	type args struct {
		key string
	}

	ims := NewInMemoryStorage()
	ims.data["m1"] = &entity.CounterMetricEntity{
		Name:  "m1",
		Value: uint64(1),
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
	ims := NewInMemoryStorage()
	ims.data["m1"] = &entity.CounterMetricEntity{
		Name:  "m1",
		Value: uint64(1),
	}
	tests := []struct {
		name string
		ims  *InMemoryStorage
		args args
		want entity.MetricEntityInterface
	}{
		{
			name: "TestInMemoryStorage_GetByKey_test1",
			ims:  ims,
			args: args{
				key: "m1",
			},
			want: ims.data["m1"],
		},
		{
			name: "TestInMemoryStorage_GetByKey_test1",
			ims:  ims,
			args: args{
				key: "m2",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ims.GetByKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryStorage.GetByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestInMemoryStorage_getAll(t *testing.T) {
// 	m1 := entity.CounterMetricEntity{
// 		Name:  "m1",
// 		Value: uint64(1),
// 	}
// 	m2 := entity.GaugeMetricEntity{
// 		Name:  "m2",
// 		Value: float64(1),
// 	}
// 	ims := NewInMemoryStorage()
// 	ims.data["m1"] = &m1
// 	ims.data["m2"] = &m2
// 	slice := make([]entity.MetricEntityInterface, len(ims.data))
// 	slice = append(slice, &m1)
// 	slice = append(slice, &m2)
// 	tests := []struct {
// 		name string
// 		ims  *InMemoryStorage
// 		want []entity.MetricEntityInterface
// 	}{
// 		{
// 			name: "TestInMemoryStorage_getAll_test1",
// 			ims:  &ims,
// 			want: slice,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.ims.getAll(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("InMemoryStorage.getAll() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
