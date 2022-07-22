package storage

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/crypto"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

func BenchmarkConvertMapsToMetrics(b *testing.B) {
	sm := NewInMemory(config.Config{})
	// sm.Lock()
	// defer sm.Unlock()

	for i := 0; i < b.N*100000; i++ {

		sm.MapCounter[strconv.Itoa(i)] = types.Counter(i)
		sm.MapGauge[strconv.Itoa(i)] = types.Gauge(i)

	}
	b.ResetTimer()
	b.Run("ConvertMapsToMetricsBeforeProfiling: Before profiling ", func(b *testing.B) {

		sm.ConvertMapsToMetricsBeforeProfiling()
	})
	b.Run("ConvertMapsToMetricsProfiled: After Profiling ", func(b *testing.B) {

		sm.ConvertMapsToMetricsProfiled()
	})
}

func (sm *InMemory) ConvertMapsToMetricsBeforeProfiling() {
	sm.Lock()
	defer sm.Unlock()

	metrics := make([]types.Metrics, len(sm.MapCounter)+len(sm.MapGauge))
	i := 0
	for k, v := range sm.MapCounter {
		vi64 := int64(v)

		metrics[i] = types.Metrics{
			ID:    k,
			MType: strings.ToLower(reflect.TypeOf(v).Name()), // исправлен после профилирования
			Delta: &vi64,
		}
		i++
	}
	for k, v := range sm.MapGauge {
		vf64 := float64(v)
		metrics[i] = types.Metrics{
			ID:    k,
			MType: strings.ToLower(reflect.TypeOf(v).Name()), // исправлен после профилирования
			Value: &vf64,
		}
		i++
	}
	// logrus.Infof("%+v", metrics)
	sm.ArrayJSONMetrics = make([]types.Metrics, len(metrics))
	copy(sm.ArrayJSONMetrics, metrics)
	logrus.Debugf("%+v", sm.ArrayJSONMetrics)
}

func (sm *InMemory) ConvertMapsToMetricsProfiled() {
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
		}
		i++
	}
	for k, v := range sm.MapGauge {
		vf64 := float64(v)
		metrics[i] = types.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &vf64,
		}
		i++
	}
	// logrus.Infof("%+v", i)
	sm.ArrayJSONMetrics = make([]types.Metrics, len(metrics))
	copy(sm.ArrayJSONMetrics, metrics)
	logrus.Debugf("%+v", sm.ArrayJSONMetrics)
}

func (sm *InMemory) UpdateJSONBeforeProfiling(cfg config.Config, metrics types.Metrics) error {
	logrus.SetReportCaller(true)

	err := crypto.CheckHash(metrics, cfg.Key)
	if err != nil {
		return fmt.Errorf("incorrect hash: %v", err)
	}

	switch metrics.MType {
	case "counter":
		if metrics.Delta == nil {
			return errors.New("recieved nil pointer on Delta")
		}
		sm.MapCounter[metrics.ID] += types.Counter(*(metrics).Delta)
		logrus.Debugf("%+v", sm.MapCounter)

	case "gauge":
		if metrics.Value == nil {
			return errors.New("recieved nil pointer on Value")
		}
		sm.MapGauge[metrics.ID] = types.Gauge(*(metrics).Value)
	}
	return nil
}

func BenchmarkCheckHash(b *testing.B) {
	cfg := config.Config{}
	sm := NewInMemory(cfg)
	metrics := types.Metrics{}
	// sm.Lock()
	// defer sm.Unlock()

	// b.ResetTimer()

	b.Run("CheckHash: Before profiling ", func(b *testing.B) {
		for i := 0; i < b.N*100000; i++ {
			sm.UpdateJSONBeforeProfiling(cfg, metrics)
		}
	})
	b.Run("CheckHash: After Profiling ", func(b *testing.B) {

		for i := 0; i < b.N*100000; i++ {
			sm.UpdateJSON(cfg, metrics)
		}

	})
}
