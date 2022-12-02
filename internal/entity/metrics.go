package entity

import (
	"strconv"
	"strings"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *uint64  `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) GetKey() string {
	return strings.ToLower(m.ID)
}

func (m Metrics) GetStringValue() string {
	var stringValue string

	switch m.MType {
	case "gauge":
		stringValue = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case "counter":
		stringValue = strconv.FormatUint(*m.Delta, 10)
	}

	return stringValue
}
