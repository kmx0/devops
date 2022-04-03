package storage

import (
	"errors"
	"strconv"

	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/types"
)

type InMemory struct {
	MapCounter map[string]types.Counter
	MapGauge   map[string]types.Gauge
}

func (s *InMemory) Update(metric string, name string, value string) error {
	switch metric {
	case "Counter":
		if _, ok := s.MapCounter[name]; ok {
			return errors.New("not such metric name")
		}
		valueInt64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		s.MapCounter[name] += types.Counter(valueInt64)
	case "Gauge":
		if _, ok := s.MapGauge[name]; ok {
			return errors.New("not such metric name")
		}
		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			// logrus.Error(err)
			return err
		}
		s.MapGauge[name] = types.Gauge(valueFloat64)
	}
	return nil
}

func NewInMemory() repositories.Repository {
	return &InMemory{
		MapCounter: make(map[string]types.Counter),
		MapGauge:   make(map[string]types.Gauge),
	}
}
