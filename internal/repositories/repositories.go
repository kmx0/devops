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
	UpdateJSON(config.Config, types.Metrics) error
	GetGauge(metric, name string) (types.Gauge, error)
	GetGaugeJSON(types.Metrics) (types.Metrics, error)
	GetCounter(metric, name string) (types.Counter, error)
	GetCounterJSON(types.Metrics) (types.Metrics, error)
	GetCurrentMetrics() (map[string]types.Gauge, map[string]types.Counter, error)
	RestoreFromDisk(cfg config.Config)
	SaveToDisk(cfg config.Config)
}
