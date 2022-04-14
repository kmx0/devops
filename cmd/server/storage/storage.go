package storage

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

type InMemory struct {
	MapCounter       map[string]types.Counter
	MapGauge         map[string]types.Gauge
	MetricNames      map[string]interface{}
	ArrayJSONMetrics []types.Metrics
	sync.RWMutex
}

func (sm *InMemory) GetCurrentMetrics() (map[string]types.Gauge, map[string]types.Counter, error) {
	return sm.MapGauge, sm.MapCounter, nil
}
func (sm *InMemory) GetGauge(metricType string, metric string) (types.Gauge, error) {
	if value, ok := sm.MapGauge[metric]; !ok {
		logrus.Info(value, ok)
		return value, errors.New("not such metric")
	}
	logrus.Info(metric)
	return sm.MapGauge[metric], nil
}
func (sm *InMemory) GetGaugeJSON(metrics types.Metrics) (types.Metrics, error) {
	logrus.Info(metrics.ID)
	if _, ok := sm.MapGauge[metrics.ID]; !ok {
		return metrics, errors.New("not such metric")
	}
	val := float64(sm.MapGauge[metrics.ID])

	metrics.Value = &val
	return metrics, nil
}
func (sm *InMemory) GetCounter(metricType string, metric string) (types.Counter, error) {
	if value, ok := sm.MapCounter[metric]; !ok {
		return value, errors.New("not such metric")
	}
	return sm.MapCounter[metric], nil
}
func (sm *InMemory) GetCounterJSON(metrics types.Metrics) (types.Metrics, error) {
	logrus.Info(metrics.ID)
	if _, ok := sm.MapCounter[metrics.ID]; !ok {
		return metrics, errors.New("not such metric")
	}
	val := int64(sm.MapCounter[metrics.ID])

	metrics.Delta = &val
	return metrics, nil
}
func (sm *InMemory) UpdateJSON(metrics types.Metrics) error {
	logrus.SetReportCaller(true)

	switch metrics.MType {
	case "counter":
		if metrics.Delta == nil {
			return errors.New("recieved nil pointer on Delta")
		}
		sm.MapCounter[metrics.ID] += types.Counter(*(metrics).Delta)
		logrus.Infof("%+v", sm.MapCounter)

	case "gauge":
		if metrics.Value == nil {
			return errors.New("recieved nil pointer on Value")
		}
		sm.MapGauge[metrics.ID] = types.Gauge(*(metrics).Value)
		// logrus.Infof("%+v", sm.MapGauge)
	}
	return nil
}
func (sm *InMemory) Update(metricType string, metric string, value string) error {
	logrus.SetReportCaller(true)

	switch metricType {
	case "counter":
		if _, ok := sm.MetricNames[metric]; !ok {
			logrus.Info("Adding new metric ", metric, " Counter")
		}
		valueInt64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		sm.MapCounter[metric] += types.Counter(valueInt64)
	case "gauge":
		if _, ok := sm.MetricNames[metric]; !ok {
			// return errors.New("not such metric")
			logrus.Info("Adding new metric ", metric, " Gauge")

		}

		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		sm.MapGauge[metric] = types.Gauge(valueFloat64)
	}
	return nil
}

func (sm *InMemory) SaveToDisk(cfg config.Config) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	sm.ConvertMapsToMetrisc()
	encoder.Encode(&sm.ArrayJSONMetrics)
	logrus.Infof("%+v", sm.ArrayJSONMetrics[0])
}
func (sm *InMemory) RestoreFromDisk(cfg config.Config) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sm.ArrayJSONMetrics)
	if err != nil {
		logrus.Error(err)
		return
	}
	sm.ConvertMetriscToMaps()
}
func NewInMemory(cfg config.Config) (*InMemory, error) {
	rm := types.RunMetrics{}
	return &InMemory{
		MapCounter:  make(map[string]types.Counter),
		MapGauge:    make(map[string]types.Gauge),
		MetricNames: rm.MapMetrics,
	}, nil
}
func (sm *InMemory) ConvertMapsToMetrisc() {
	sm.Lock()
	defer sm.Unlock()
	metrics := make([]types.Metrics, len(sm.MapCounter)+len(sm.MapGauge))
	// val := reflect.ValueOf(rm)
	i := 0
	for k, v := range sm.MapCounter {
		vi64 := int64(v)

		metrics[i] = types.Metrics{
			ID:    k,
			MType: strings.ToLower(reflect.TypeOf(v).Name()),
			Delta: &vi64,
		}
		// vf64, ok := v.(float64)

		i++
	}
	for k, v := range sm.MapGauge {
		vf64 := float64(v)
		metrics[i] = types.Metrics{
			ID:    k,
			MType: strings.ToLower(reflect.TypeOf(v).Name()),
			Value: &vf64,
		}
		i++
	}
	copy(sm.ArrayJSONMetrics, metrics)
	
	// return metrics
}

func (sm *InMemory) ConvertMetriscToMaps() {
	sm.Lock()
	defer sm.Unlock()
	for _, v := range sm.ArrayJSONMetrics {

		logrus.Info(v)

		switch v.MType {
		case "counter":
			sm.MapCounter[v.ID] = types.Counter(*v.Delta)
		case "gauge":
			sm.MapGauge[v.ID] = types.Gauge(*v.Value)
		}

	}
}
