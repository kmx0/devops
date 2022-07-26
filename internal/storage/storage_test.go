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
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
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
			_ = sm.UpdateJSONBeforeProfiling(cfg, metrics)
		}
	})
	b.Run("CheckHash: After Profiling ", func(b *testing.B) {

		for i := 0; i < b.N*100000; i++ {
			_ = sm.UpdateJSON(cfg.Key, metrics)
		}

	})
}

func TestGetCurrentMetrics(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapGauge["1"] = types.Gauge(1)
	s.MapCounter["1"] = types.Counter(1)
	wantC := make(map[string]types.Counter)
	wantG := make(map[string]types.Gauge)
	wantC["1"] = types.Counter(1)
	wantG["1"] = types.Gauge(1)
	type wantStruct struct {
		gk string
		ck string
		gv types.Gauge
		cv types.Counter
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Simple test 1",
			want: wantStruct{
				gk: "1",
				ck: "1",
				gv: types.Gauge(1),
				cv: types.Counter(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resG, resC, err := s.GetCurrentMetrics()

			assert.Equal(t, tt.want.gv, resG[tt.want.gk])
			assert.Equal(t, tt.want.cv, resC[tt.want.ck])
			require.NoError(t, err)
		})
	}
}

func TestGetGauge(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapGauge["Alloc"] = types.Gauge(1)
	type wantStruct struct {
		gv  types.Gauge
		err error
	}

	tests := []struct {
		name   string
		metric string
		want   wantStruct
	}{
		{
			name:   "Alloc",
			metric: "Alloc",
			want: wantStruct{
				gv:  types.Gauge(1),
				err: nil,
			},
		},
		{
			name:   "FailAlloc",
			metric: "FailAlloc",
			want: wantStruct{
				err: errors.New("not such metric"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resG, err := s.GetGauge("gauge", tt.metric)

			assert.Equal(t, tt.want.gv, resG)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetGaugeJSON(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapGauge["Alloc"] = types.Gauge(1)
	type wantStruct struct {
		value float64
		err   error
	}
	tests := []struct {
		name     string
		metricID string
		want     wantStruct
	}{
		{
			name: "Alloc",

			metricID: "Alloc",
			want: wantStruct{
				value: float64(1),
				err:   nil,
			},
		},
		{
			name:     "FailAlloc",
			metricID: "FailAlloc",
			want: wantStruct{
				err: errors.New("not such metric"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := s.GetGaugeJSON(tt.metricID)

			assert.Equal(t, tt.want.value, value)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetCounter(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapCounter["PollCount"] = types.Counter(1)
	type wantStruct struct {
		cv  types.Counter
		err error
	}

	tests := []struct {
		name   string
		metric string
		want   wantStruct
	}{
		{
			name:   "PollCount",
			metric: "PollCount",
			want: wantStruct{
				cv:  types.Counter(1),
				err: nil,
			},
		},
		{
			name:   "FailPollCount",
			metric: "FailPollCount",
			want: wantStruct{
				err: errors.New("not such metric"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resC, err := s.GetCounter("counter", tt.metric)

			assert.Equal(t, tt.want.cv, resC)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetCounterJSON(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapCounter["PollCount"] = types.Counter(1)
	type wantStruct struct {
		delta int64
		err   error
	}
	tests := []struct {
		name     string
		metricID string
		want     wantStruct
	}{
		{
			name: "PollCount",

			metricID: "PollCount",
			want: wantStruct{
				delta: int64(1),
				err:   nil,
			},
		},
		{
			name:     "FailPollCont",
			metricID: "FailPollCount",
			want: wantStruct{
				err: errors.New("not such metric"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta, err := s.GetCounterJSON(tt.metricID)

			assert.Equal(t, tt.want.delta, delta)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateJSON(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapCounter["PollCount"] = types.Counter(1)
	type wantStruct struct {
		err error
	}
	var helperf float64 = 1
	var helperi int64 = 1

	tests := []struct {
		name    string
		hashkey string
		metrics types.Metrics
		want    wantStruct
	}{
		{
			name: "Correct test Gauge",

			hashkey: "hashkey",
			metrics: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: &helperf,
				Hash:  crypto.Hash(fmt.Sprintf("%s:gauge:%f", "Alloc", helperf), "hashkey"),
			},
			want: wantStruct{
				err: nil,
			},
		},
		{
			name: "Correct test Counter",

			hashkey: "hashkey",
			metrics: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: &helperi,
				Hash:  crypto.Hash(fmt.Sprintf("%s:counter:%d", "PollCount", helperi), "hashkey"),
			},
			want: wantStruct{
				err: nil,
			},
		},
		{
			name: "Incorrect Hash",

			hashkey: "hashkey",
			metrics: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: &helperf,
				Hash:  crypto.Hash(fmt.Sprintf("%s:gauge:%f", "Alloc", helperf), "Failhashkey"),
			},
			want: wantStruct{
				err: fmt.Errorf("incorrect hash: %v", errors.New("hash sum not matched")),
			},
		},
		{
			name: "Nil Pointer on Value",

			metrics: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
			},
			want: wantStruct{
				err: errors.New("recieved nil pointer on Value"),
			},
		},
		{
			name: "Nil Pointer on Delta",

			metrics: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
			},
			want: wantStruct{
				err: errors.New("recieved nil pointer on Delta"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateJSON(tt.hashkey, tt.metrics)

			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapCounter["PollCount"] = types.Counter(1)
	type wantStruct struct {
		err error
	}

	tests := []struct {
		name       string
		metricType string
		metric     string
		value      string
		want       wantStruct
	}{
		{
			name: "Correct test Gauge",

			metricType: "gauge",
			metric:     "Alloc",
			value:      "1.0",
			want: wantStruct{
				err: nil,
			},
		},
		{
			name: "Correct test Counter",

			metricType: "counter",
			metric:     "PollCount",
			value:      "1",
			want: wantStruct{
				err: nil,
			},
		},
		{
			name: "Incorrect test Gauge",

			metricType: "gauge",
			metric:     "Alloc",
			value:      "fail",
			want: wantStruct{
				err: errors.New("strconv.ParseFloat: parsing"),
			},
		},
		{
			name: "Incorrect test Counter",

			metricType: "counter",
			metric:     "PollCount",
			value:      "fail",
			want: wantStruct{
				err: errors.New("strconv.ParseInt: parsing"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// metricType string, metric string, value string
			err := s.Update(tt.metricType, tt.metric, tt.value)

			if err != nil {
				assert.Equal(t, strings.Contains(err.Error(), tt.want.err.Error()), true)

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConvertMapsToMetrics(t *testing.T) {
	s := NewInMemory(config.Config{})
	s.MapGauge["1"] = types.Gauge(1)
	s.MapCounter["1"] = types.Counter(1)
	type wantStruct struct {
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Simple test 1",
			want: wantStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ConvertMapsToMetrics()

			assert.Equal(t, len(s.ArrayJSONMetrics), len(s.MapCounter)+len(s.MapGauge))
		})
	}
}

func TestConvertMetricsToMaps(t *testing.T) {
	s := NewInMemory(config.Config{})
	// s.MapGauge["1"] = types.Gauge(1)
	// s.MapCounter["1"] = types.Counter(1)
	var helperf float64 = 1
	var helperi int64 = 1
	s.ArrayJSONMetrics = append(s.ArrayJSONMetrics, types.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &helperf,
	})

	s.ArrayJSONMetrics = append(s.ArrayJSONMetrics, types.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &helperi,
	})
	type wantStruct struct {
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Simple test 1",
			want: wantStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.ConvertMetricsToMaps()

			assert.Equal(t, len(s.ArrayJSONMetrics), len(s.MapCounter)+len(s.MapGauge))
		})
	}
}
