package storage

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

type InMemory struct {
	MapCounter       map[string]types.Counter
	MapGauge         map[string]types.Gauge
	MetricNames      map[string]interface{}
	ArrayJSONMetrics []types.Metrics
	file             *os.File
	encoder          *json.Encoder
	sync.RWMutex
}
type InDisk struct {
	MapCounter  map[string]types.Counter
	MapGauge    map[string]types.Gauge
	MetricNames map[string]interface{}
	file        *os.File
	decoder     *json.Decoder
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
		logrus.Infof("%+v", sm.MapGauge)
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

func (sm *InMemory) WriteMetrics(metricsP *[]types.Metrics) error {
	err := sm.encoder.Encode(metricsP)
	if err != nil {
		return err
	}
	return nil
}

func NewInMemory(filename string) (*InMemory, error) {
	rm := types.RunMetrics{}
	// val := reflect.ValueOf(rm)
	// metricNames := make(map[string]struct{}, val.NumField())
	// // val := reflect.ValueOf(rm)
	// for i := 0; i < val.NumField(); i++ {
	// 	metricNames[val.Type().Field(i).Name] = struct{}{}
	// }
	if filename == "" {
		return &InMemory{
			MapCounter:  make(map[string]types.Counter),
			MapGauge:    make(map[string]types.Gauge),
			MetricNames: rm.MapMetrics,
		}, nil
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &InMemory{
		MapCounter:  make(map[string]types.Counter),
		MapGauge:    make(map[string]types.Gauge),
		MetricNames: rm.MapMetrics,
		file:        file,
		encoder:     json.NewEncoder(file),
	}, nil

}

func ConvertMapsToMetrisc(sm *InMemory) []types.Metrics {
	metrics := make([]types.Metrics, len(sm.MapCounter)+len(sm.MapGauge))
	sm.Lock()
	defer sm.Unlock()
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
	return metrics
}

func (sm *InMemory) Close() error {
	return sm.file.Close()
}

func (sd *InDisk) ReadMetrics() (*[]types.Metrics, error) {
	metrics := []types.Metrics{}
	err := sd.decoder.Decode(&metrics)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &metrics, nil
}
func (sd *InDisk) Close() error {
	return sd.file.Close()
}

func NewInDisk(filename string) (*InDisk, error) {
	rm := types.RunMetrics{}
	// val := reflect.ValueOf(rm)
	// metricNames := make(map[string]struct{}, val.NumField())
	// // val := reflect.ValueOf(rm)
	// for i := 0; i < val.NumField(); i++ {
	// 	metricNames[val.Type().Field(i).Name] = struct{}{}
	// }
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &InDisk{
		MapCounter:  make(map[string]types.Counter),
		MapGauge:    make(map[string]types.Gauge),
		MetricNames: rm.MapMetrics,
		file:        file,
		decoder:     json.NewDecoder(file),
	}, nil
}
