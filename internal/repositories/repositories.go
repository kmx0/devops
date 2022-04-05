package repositories

import "github.com/kmx0/devops/internal/types"

type Repository interface {
	Update(metric, name, value string) error
	GetGauge(metric, name string) (types.Gauge, error)
	GetCounter(metric, name string) (types.Counter, error)
	GetCurrentMetrics() (map[string]types.Gauge, map[string]types.Counter, error)
}
