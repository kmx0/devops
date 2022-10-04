package types

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

type (
	Gauge float64

	Counter int64
	// RunMetrics - strcut for filling metrics from Memory Stats
	RunMetrics struct {
		Alloc           Gauge
		BuckHashSys     Gauge
		Frees           Gauge
		GCCPUFraction   Gauge
		GCSys           Gauge
		HeapAlloc       Gauge
		HeapIdle        Gauge
		HeapInuse       Gauge
		HeapObjects     Gauge
		HeapReleased    Gauge
		HeapSys         Gauge
		LastGC          Gauge
		Lookups         Gauge
		MCacheInuse     Gauge
		MCacheSys       Gauge
		MSpanInuse      Gauge
		MSpanSys        Gauge
		Mallocs         Gauge
		NextGC          Gauge
		NumForcedGC     Gauge
		NumGC           Gauge
		OtherSys        Gauge
		PauseTotalNs    Gauge
		StackInuse      Gauge
		StackSys        Gauge
		Sys             Gauge
		TotalAllo       Gauge
		PollCount       Counter
		RandomValue     Gauge
		TotalMemory     Gauge
		FreeMemory      Gauge
		CPUutilization1 Gauge

		sync.RWMutex
		MapMetrics map[string]interface{}
	}
)

// Metrics - struct for sending data use JSON format
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (g Gauge) String() string {
	return fmt.Sprintf("%.3f", g)
}
func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

// GetMetrics - convert metrics from RunMetrics struct to Metrics struct for sending use JSON format
func (rm *RunMetrics) GetMetrics() (metricsForBody []Metrics) {

	metrics := make([]Metrics, len(rm.MapMetrics))
	rm.Lock()
	defer rm.Unlock()
	i := 0
	for k, v := range rm.MapMetrics {
		ty := strings.ToLower(reflect.TypeOf(v).Name()) // тип метрики
		switch ty {
		case "counter":
			vc, ok := v.(Counter)
			if !ok {

				logrus.Error(errors.New("cannot convert interface to int64"))
			}
			vi64 := int64(vc)

			metrics[i] = Metrics{
				ID:    k,
				MType: "counter",
				Delta: &vi64,
			}
		case "gauge":
			vg, ok := v.(Gauge)
			if !ok {

				logrus.Error(errors.New("cannot convert interface to float64"))
			}
			vf64 := float64(vg)
			metrics[i] = Metrics{
				ID:    k,
				MType: "gauge",
				Value: &vf64,
			}

		}
		i++
	}
	return metrics
}

// Set - function for setting metrics from Memory Stats
// Use in agent code
func (rm *RunMetrics) Set(ms runtime.MemStats) {
	rm.Lock()
	defer rm.Unlock()
	rm.MapMetrics["Alloc"] = Gauge(ms.Alloc)
	rm.MapMetrics["BuckHashSys"] = Gauge(ms.BuckHashSys)
	rm.MapMetrics["Frees"] = Gauge(ms.Frees)
	rm.MapMetrics["GCCPUFraction"] = Gauge(ms.GCCPUFraction)
	rm.MapMetrics["GCSys"] = Gauge(ms.GCSys)
	rm.MapMetrics["HeapAlloc"] = Gauge(ms.HeapAlloc)
	rm.MapMetrics["HeapIdle"] = Gauge(ms.HeapIdle)
	rm.MapMetrics["HeapInuse"] = Gauge(ms.HeapInuse)
	rm.MapMetrics["HeapObjects"] = Gauge(ms.HeapObjects)
	rm.MapMetrics["HeapReleased"] = Gauge(ms.HeapReleased)
	rm.MapMetrics["HeapSys"] = Gauge(ms.HeapSys)
	rm.MapMetrics["LastGC"] = Gauge(ms.LastGC)
	rm.MapMetrics["Lookups"] = Gauge(ms.Lookups)
	rm.MapMetrics["MCacheInuse"] = Gauge(ms.MCacheInuse)
	rm.MapMetrics["MCacheSys"] = Gauge(ms.MCacheSys)
	rm.MapMetrics["MSpanInuse"] = Gauge(ms.MSpanInuse)
	rm.MapMetrics["MSpanSys"] = Gauge(ms.MSpanSys)
	rm.MapMetrics["Mallocs"] = Gauge(ms.Mallocs)
	rm.MapMetrics["NextGC"] = Gauge(ms.NextGC)
	rm.MapMetrics["NumForcedGC"] = Gauge(ms.NumForcedGC)
	rm.MapMetrics["NumGC"] = Gauge(ms.NumGC)
	rm.MapMetrics["OtherSys"] = Gauge(ms.OtherSys)
	rm.MapMetrics["PauseTotalNs"] = Gauge(ms.PauseTotalNs)
	rm.MapMetrics["StackInuse"] = Gauge(ms.StackInuse)
	rm.MapMetrics["StackSys"] = Gauge(ms.StackSys)
	rm.MapMetrics["Sys"] = Gauge(ms.Sys)
	rm.MapMetrics["TotalAlloc"] = Gauge(ms.TotalAlloc)
	if rm.MapMetrics["PollCount"] == nil {
		rm.MapMetrics["PollCount"] = Counter(0)
	}
	rm.MapMetrics["PollCount"] = (rm.MapMetrics["PollCount"].(Counter)) + Counter(1)
	rand.Seed(time.Now().UnixNano())
	rm.MapMetrics["RandomValue"] = Gauge(rand.Float64())
}

// Set - function for setting metrics from Gopsitil
// Use in agent code
func (rm *RunMetrics) SetGopsutil() {
	rm.Lock()
	defer rm.Unlock()
	v, _ := mem.VirtualMemory()
	p, _ := cpu.Percent(time.Second, false)
	rm.MapMetrics["TotalMemory"] = Gauge(v.Total)
	rm.MapMetrics["FreeMemory"] = Gauge(v.Free)
	rm.MapMetrics["CPUutilization1"] = Gauge(p[0])
}
