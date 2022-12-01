package entity

import (
	"fmt"
	"strconv"
	"strings"
)

type GaugeMetricEntity struct {
	Name  string
	Value float64
}

func (gme *GaugeMetricEntity) GetUpdateURI() string {
	return fmt.Sprintf("/update/gauge/%s/%f", gme.Name, gme.Value)
}

func (gme *GaugeMetricEntity) GetValue() interface{} {
	return gme.Value
}

func (gme *GaugeMetricEntity) GetKey() string {
	return strings.ToLower(gme.Name)
}

func (gme *GaugeMetricEntity) SetValue(value interface{}) {
	floatValue, ok := value.(float64)
	if ok {
		gme.Value = floatValue
	}
}

func (gme *GaugeMetricEntity) GetStringValue() string {
	return strconv.FormatFloat(gme.Value, 'f', -1, 64)
}

func (gme *GaugeMetricEntity) GetShortTypeName() string {
	return "gauge"
}

func (gme *GaugeMetricEntity) GetMetricEntity() Metrics {
	return Metrics{
		ID:    gme.Name,
		MType: gme.GetShortTypeName(),
		Value: &gme.Value,
	}
}

func (gme *GaugeMetricEntity) GetName() string {
	return gme.Name
}
