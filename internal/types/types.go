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

	"github.com/sirupsen/logrus"
)

type (
	Gauge float64

	Counter int64

	RunMetrics struct {
		Alloc         Gauge
		BuckHashSys   Gauge
		Frees         Gauge
		GCCPUFraction Gauge
		GCSys         Gauge
		HeapAlloc     Gauge
		HeapIdle      Gauge
		HeapInuse     Gauge
		HeapObjects   Gauge
		HeapReleased  Gauge
		HeapSys       Gauge
		LastGC        Gauge
		Lookups       Gauge
		MCacheInuse   Gauge
		MCacheSys     Gauge
		MSpanInuse    Gauge
		MSpanSys      Gauge
		Mallocs       Gauge
		NextGC        Gauge
		NumForcedGC   Gauge
		NumGC         Gauge
		OtherSys      Gauge
		PauseTotalNs  Gauge
		StackInuse    Gauge
		StackSys      Gauge
		Sys           Gauge
		TotalAlloc    Gauge
		PollCount     Counter
		RandomValue   Gauge
		sync.RWMutex
		MapMetrics map[string]interface{}
	}
)
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (g Gauge) String() string {
	return fmt.Sprintf("%.3f", g)
}
func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

func (rm *RunMetrics) Get() (endpoint string, metricsForBody []Metrics) {

	metrics := make([]Metrics, len(rm.MapMetrics))
	rm.Lock()
	defer rm.Unlock()
	// val := reflect.ValueOf(rm)
	i := 0
	for k, v := range rm.MapMetrics {
		// endpoint :=

		// endpoint = fmt.Sprintf("%s/%s/%s/%v", endpoint, strings.ToLower(reflect.TypeOf(v).Name()), k, v)
		//counter or gauge

		ty := strings.ToLower(reflect.TypeOf(v).Name()) //тип метрики
		switch ty {
		case "counter":
			vc, ok := v.(Counter)
			if !ok {

				panic(errors.New("cannot convert interface to int64"))
			}
			vi64 := int64(vc)
			logrus.Errorf("%+v", v)
			logrus.Errorf("%+v", vi64)
			metrics[i] = Metrics{
				ID:    k,
				MType: strings.ToLower(reflect.TypeOf(v).Name()),
				Delta: &vi64,
			}
		case "gauge":
			// vf64, ok := v.(float64)
			vg, ok := v.(Gauge)
			if !ok {

				panic(errors.New("cannot convert interface to float64"))
			}
			vf64 := float64(vg)
			// logrus.Errorf("%+v", v)
			// logrus.Errorf("%+v", vf64)
			metrics[i] = Metrics{
				ID:    k,
				MType: strings.ToLower(reflect.TypeOf(v).Name()),
				Value: &vf64,
			}

		}
		i++
	}

	// endpoint = fmt.Sprintf("%s/%s/%s/%v", endpoint, strings.ToLower(val.Type().Field(i).Type.Name()), val.Type().Field(i).Name, rmMap[val.Type().Field(i).Name])

	// АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	return "http://127.0.0.1:8080/update", metrics
}

func (rm *RunMetrics) Set(ms runtime.MemStats) {
	rm.Lock()
	defer rm.Unlock()
	// rm.MapMetrics["MapMetrics["Alloc"] = Gauge(ms.Alloc)
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
	rm.RandomValue = Gauge(rand.Float64())

}
