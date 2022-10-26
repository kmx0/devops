package repositories

import (
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/types"
)

// Repository - interface for storing data
// Exist two struct for implementing its interface
// File or Postgres DB
type Repository interface {
	Update(metric, name, value string) error
	UpdateJSON(string, types.Metrics) error
	GetGauge(metric, name string) (types.Gauge, error)
	GetGaugeJSON(string) (float64, error)
	GetCounter(metric, name string) (types.Counter, error)
	GetCounterJSON(string) (int64, error)
	GetCurrentMetrics() (map[string]types.Gauge, map[string]types.Counter, error)
	RestoreFromDisk(cfg config.Config) error
	SaveToDisk(cfg config.Config) error
}
