package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/crypto"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

// InMemory - implements Repository interface
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
		logrus.Debug(value, ok)
		return value, errors.New("not such metric")
	}
	return sm.MapGauge[metric], nil
}
func (sm *InMemory) GetGaugeJSON(metricID string) (float64, error) {
	if _, ok := sm.MapGauge[metricID]; !ok {
		return 0, errors.New("not such metric")
	}
	val := float64(sm.MapGauge[metricID])
	return val, nil
}

func (sm *InMemory) GetCounter(metricType string, metric string) (types.Counter, error) {
	if value, ok := sm.MapCounter[metric]; !ok {
		return value, errors.New("not such metric")
	}

	return sm.MapCounter[metric], nil
}
func (sm *InMemory) GetCounterJSON(metricID string) (int64, error) {
	if _, ok := sm.MapCounter[metricID]; !ok {
		return 0, errors.New("not such metric")
	}
	delta := int64(sm.MapCounter[metricID])

	return delta, nil
}

// UpdateJSON - check hash in JSON
// saving metrics to Maps
func (sm *InMemory) UpdateJSON(hashkey string, metrics types.Metrics) error {
	logrus.SetReportCaller(true)

	if hashkey != "" {
		err := crypto.CheckHash(metrics, hashkey)
		if err != nil {
			return fmt.Errorf("incorrect hash: %w", err)
		}
	}
	switch metrics.MType {
	case "counter":
		if metrics.Delta == nil {
			return errors.New("received nil pointer on Delta")
		}
		sm.MapCounter[metrics.ID] += types.Counter(*(metrics).Delta)
		logrus.Debugf("%+v", sm.MapCounter)

	case "gauge":
		if metrics.Value == nil {
			return errors.New("received nil pointer on Value")
		}
		sm.MapGauge[metrics.ID] = types.Gauge(*(metrics).Value)
	}
	return nil
}

// Update - saving metrics to Maps without checking hash
// not using JSON struct
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

// SaveToDisk - saving metrics from Maps to file or DB
func (sm *InMemory) SaveToDisk(cfg config.Config) error {
	if cfg.DBDSN == "" {
		file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			logrus.Error(err)
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		sm.ConvertMapsToMetrics()

		if err = encoder.Encode(&sm.ArrayJSONMetrics); err != nil {
			logrus.Errorf("error saving to disk: %v", err)
		}
	}
	if cfg.DBDSN != "" {
		//saving to db
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		PingDB(ctx, cfg.DBDSN)
		err := SaveDataToDB(sm)
		if err != nil {
			return err
		}
		logrus.Info("Saving to DB")

	}
	return nil
}

// RestoreFromDisk - Get Metrics from storage before start server, if flag Restore = true
func (sm *InMemory) RestoreFromDisk(cfg config.Config) error {
	if cfg.DBDSN == "" {
		file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			logrus.Error(err)
			return err
		}
		defer file.Close()
		decoder := json.NewDecoder(file)

		if decoder == nil {
			return errors.New("file is empty")
		}
		err = decoder.Decode(&sm.ArrayJSONMetrics)
		if err != nil {
			logrus.Error(err)
			return err
		}
		sm.ConvertMetricsToMaps()
	}
	if cfg.DBDSN != "" {
		//saving to db
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// не забываем освободить ресурс
		defer cancel()
		if PingDB(ctx, cfg.DBDSN) {
			err := RestoreDataFromDB(sm)
			if err != nil {
				return err
			}
			logrus.Info("Restoring from DB")
		} else {
			return errors.New("can not connect to DB")
		}
	}
	return nil
}

// Constructor for InMemory
func NewInMemory(cfg config.Config) *InMemory {
	rm := types.RunMetrics{}
	return &InMemory{
		MapCounter:       make(map[string]types.Counter),
		MapGauge:         make(map[string]types.Gauge),
		MetricNames:      rm.MapMetrics,
		ArrayJSONMetrics: make([]types.Metrics, 0),
	}
}

// ConvertMapsToMetrics - Converting from Maps to Metrics struct for JSON format
func (sm *InMemory) ConvertMapsToMetrics() {
	sm.Lock()
	defer sm.Unlock()

	metrics := make([]types.Metrics, len(sm.MapCounter)+len(sm.MapGauge))
	i := 0
	for k, v := range sm.MapCounter {
		vi64 := int64(v)

		metrics[i] = types.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &vi64,
			Value: new(float64),
			Hash:  "",
		}
		i++
	}
	for k, v := range sm.MapGauge {
		vf64 := float64(v)
		metrics[i] = types.Metrics{
			ID: k,
			// MType: strings.ToLower(reflect.TypeOf(v).Name()),// исправлен после профилирования
			MType: "gauge",
			Value: &vf64,
		}
		i++
	}
	sm.ArrayJSONMetrics = make([]types.Metrics, len(metrics))
	copy(sm.ArrayJSONMetrics, metrics)
	logrus.Debugf("%+v", sm.ArrayJSONMetrics)
}

// ConvertMetricsToMaps - Converting from Metrics struct to Maps
func (sm *InMemory) ConvertMetricsToMaps() {
	sm.Lock()
	defer sm.Unlock()
	for _, v := range sm.ArrayJSONMetrics {
		switch v.MType {
		case "counter":
			sm.MapCounter[v.ID] = types.Counter(*v.Delta)
		case "gauge":
			sm.MapGauge[v.ID] = types.Gauge(*v.Value)
		}

	}
}
