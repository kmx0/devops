package storage

import (
	"errors"
	"strconv"

	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

type InMemory struct {
	MapCounter  map[string]types.Counter
	MapGauge    map[string]types.Gauge
	MetricNames map[string]interface{}
}

func (s *InMemory) GetCurrentMetrics() (map[string]types.Gauge, map[string]types.Counter, error) {
	return s.MapGauge, s.MapCounter, nil
}
func (s *InMemory) GetGauge(metricType string, metric string) (types.Gauge, error) {
	if value, ok := s.MapGauge[metric]; !ok {
		logrus.Info(value, ok)
		return value, errors.New("not such metric")
	}
	logrus.Info(metric)
	return s.MapGauge[metric], nil
}
func (s *InMemory) GetGaugeJSON(metrics types.Metrics) (types.Metrics, error) {
	if _, ok := s.MapCounter[metrics.ID]; !ok {
		return metrics, errors.New("not such metric")
	}
	val := float64(s.MapCounter[metrics.ID])

	metrics.Value = &val
	// return s.MapCounter[metric], nil
	return metrics, nil
}
func (s *InMemory) GetCounter(metricType string, metric string) (types.Counter, error) {
	if value, ok := s.MapCounter[metric]; !ok {
		return value, errors.New("not such metric")
	}
	return s.MapCounter[metric], nil
}
func (s *InMemory) GetCounterJSON(metrics types.Metrics) (types.Metrics, error) {
	if _, ok := s.MapCounter[metrics.ID]; !ok {
		return metrics, errors.New("not such metric")
	}
	val := int64(s.MapCounter[metrics.ID])

	metrics.Delta = &val
	// return s.MapCounter[metric], nil
	return metrics, nil
}
func (s *InMemory) UpdateJSON(metrics types.Metrics) error {
	logrus.SetReportCaller(true)

	switch metrics.MType {
	case "counter":
		if _, ok := s.MetricNames[metrics.ID]; !ok {
			logrus.Info("Adding new metric ", metrics.ID, " Counter")
		}
		// c := metrics.Delta

		// valueInt64, err := strconv.ParseInt(*c, 10, 64)
		// if err != nil {
		// 	return err
		// }
		logrus.Warn(*(metrics).Delta)
		s.MapCounter[metrics.ID] += types.Counter(*(metrics).Delta)
	case "gauge":
		if _, ok := s.MetricNames[metrics.ID]; !ok {
			// return errors.New("not such metric")
			logrus.Info("Adding new metric ", metrics.ID, " Gauge")

		}

		// valueFloat64, err := strconv.ParseFloat(value, 64)
		// if err != nil {
		// 	return err
		// }
		s.MapGauge[metrics.ID] = types.Gauge(*(metrics).Value)
	}
	return nil
}
func (s *InMemory) Update(metricType string, metric string, value string) error {
	logrus.SetReportCaller(true)

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
			// return errors.New("not such metric")
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
	rm := types.RunMetrics{}
	// val := reflect.ValueOf(rm)
	// metricNames := make(map[string]struct{}, val.NumField())
	// // val := reflect.ValueOf(rm)
	// for i := 0; i < val.NumField(); i++ {
	// 	metricNames[val.Type().Field(i).Name] = struct{}{}
	// }
	return &InMemory{
		MapCounter:  make(map[string]types.Counter),
		MapGauge:    make(map[string]types.Gauge),
		MetricNames: rm.MapMetrics,
	}
}
