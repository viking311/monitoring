package entity

import (
	"fmt"
	"strings"
)

type CounterMetricEntity struct {
	Name  string
	Value uint64
}

func (cme *CounterMetricEntity) GetUpdateURI() string {
	return fmt.Sprintf("/update/counter/%s/%d", cme.Name, cme.Value)
}

func (cme *CounterMetricEntity) GetValue() interface{} {
	return cme.Value
}

func (cme *CounterMetricEntity) SetValue(value interface{}) {
	intValue, ok := value.(uint64)
	if ok {
		cme.Value += intValue
	}
}

func (cme *CounterMetricEntity) GetKey() string {
	return strings.ToLower(cme.Name)
}

func (cme *CounterMetricEntity) GetStringValue() string {
	return fmt.Sprintf("%d", cme.Value)
}

func (cme *CounterMetricEntity) GetShortTypeName() string {
	return "counter"
}

func (cme *CounterMetricEntity) GetMetricEntity() Metrics {
	return Metrics{
		ID:    cme.Name,
		MType: cme.GetShortTypeName(),
		Delta: &cme.Value,
	}
}
