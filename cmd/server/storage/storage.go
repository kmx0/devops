package storage

import (
	"strconv"

	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

type InMemory struct {
	MapCounter  map[string]types.Counter
	MapGauge    map[string]types.Gauge
	MetricNames map[string]struct{}
}

func (s *InMemory) Update(metricType string, metric string, value string) error {
	switch metricType {
	case "counter":
		if _, ok := s.MetricNames[metric]; !ok {
			logrus.Info("Adding new metric ", metric, " Counter")
		}
		valueInt64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		s.MapCounter[metric] += types.Counter(valueInt64)
	case "gauge":
		if _, ok := s.MetricNames[metric]; !ok {
			logrus.Info("Adding new metric ", metric, " Gauge")
		}

		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		s.MapGauge[metric] = types.Gauge(valueFloat64)
	}
	return nil
}

func NewInMemory() repositories.Repository {
	// rm := types.RunMetrics{}
	// val := reflect.ValueOf(rm)
	// metricNames := make(map[string]struct{}, val.NumField())
	// // val := reflect.ValueOf(rm)
	// for i := 0; i < val.NumField(); i++ {
	// 	metricNames[val.Type().Field(i).Name] = struct{}{}
	// }
	return &InMemory{
		MapCounter: make(map[string]types.Counter),
		MapGauge:   make(map[string]types.Gauge),
		// MetricNames: metricNames,
	}
}
