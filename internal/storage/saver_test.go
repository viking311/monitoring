package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viking311/monitoring/internal/entity"
)

func TestSaver_Go(t *testing.T) {
	ch := make(chan entity.MetricEntityInterface, 100)
	type args struct {
		value entity.MetricEntityInterface
	}
	m1 := &entity.CounterMetricEntity{
		Name:  "m1",
		Value: uint64(1),
	}
	tests := []struct {
		name string
		s    Saver
		args args
		want entity.MetricEntityInterface
	}{
		{
			name: "TestSaver_Go_test1",
			s:    NewSaver(ch),
			args: args{
				value: m1,
			},
			want: m1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.s.Go()
			ch <- tt.args.value
			time.Sleep(1 * time.Second)
			val := tt.s.storage.GetByKey("m1")
			assert.Equal(t, tt.want, val)
		})
	}
}
