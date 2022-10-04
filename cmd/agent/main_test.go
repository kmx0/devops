package main

import (
	"errors"
	"runtime"
	"testing"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSendMetricsJSON(t *testing.T) {
	logrus.SetReportCaller(true)
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	rm.Set(m)
	rm.SetGopsutil()

	type wantStruct struct {
		err error
	}
	var helperi int64 = 1
	var helperf float64 = 1
	tests := []struct {
		cfg     config.Config
		metricp types.Metrics
		name    string
		want    wantStruct
	}{
		{
			name: "Fail",
			cfg:  config.Config{Key: "hashkey"},
			metricp: types.Metrics{
				MType: "fail",
			},
			want: wantStruct{
				err: errors.New("unknown metric type"),
			},
		},
		{
			name: "Counter",
			cfg:  config.Config{Key: "hashkey"},
			metricp: types.Metrics{
				MType: "counter",
				Delta: &helperi,
			},
			want: wantStruct{
				err: nil,
			},
		},
		{
			name: "Gauge",
			cfg:  config.Config{Key: "hashkey"},
			metricp: types.Metrics{
				MType: "gauge",
				Value: &helperf,
			},
			want: wantStruct{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendMetricsJSON(tt.cfg)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
		})
	}
}

func TestAddHash(t *testing.T) {

	key := ""

	type wantStruct struct {
		err error
	}
	var helperi int64 = 1
	var helperf float64 = 1
	tests := []struct {
		metricp types.Metrics
		name    string
		want    wantStruct
	}{
		{
			name: "Fail",
			metricp: types.Metrics{
				MType: "fail",
			},
			want: wantStruct{
				err: errors.New("unknown metric type"),
			},
		},
		{
			name: "Counter",
			metricp: types.Metrics{
				MType: "counter",
				Delta: &helperi,
			},
			want: wantStruct{
				err: nil,
			},
		},
		{
			name: "Gauge",
			metricp: types.Metrics{
				MType: "gauge",
				Value: &helperf,
			},
			want: wantStruct{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddHash(key, &tt.metricp)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
		})
	}
}
