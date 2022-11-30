package entity

type MetricEntityInterface interface {
	GetUpdateURI() string
	GetValue() interface{}
	SetValue(value interface{})
	GetKey() string
	GetStringValue() string
	GetShortTypeName() string
	GetMetricEntity() Metrics
}
