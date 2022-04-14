package repositories

import "github.com/kmx0/devops/internal/types"

type Repository interface {
	Update(metric, name, value string) error
	UpdateJSON(types.Metrics) error
	GetGauge(metric, name string) (types.Gauge, error)
	GetGaugeJSON(types.Metrics) (types.Metrics, error)
	GetCounter(metric, name string) (types.Counter, error)
	GetCounterJSON(types.Metrics) (types.Metrics, error)
	GetCurrentMetrics() (map[string]types.Gauge, map[string]types.Counter, error)
}
