package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/structs"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

func fill(ms runtime.MemStats, rm *types.RunMetrics) {
	rm.Alloc = types.Gauge(ms.Alloc)
	rm.BuckHashSys = types.Gauge(ms.BuckHashSys)
	rm.Frees = types.Gauge(ms.Frees)
	rm.GCCPUFraction = types.Gauge(ms.GCCPUFraction)
	rm.GCSys = types.Gauge(ms.GCSys)
	rm.HeapAlloc = types.Gauge(ms.HeapAlloc)
	rm.HeapIdle = types.Gauge(ms.HeapIdle)
	rm.HeapInuse = types.Gauge(ms.HeapInuse)
	rm.HeapObjects = types.Gauge(ms.HeapObjects)
	rm.HeapReleased = types.Gauge(ms.HeapReleased)
	rm.HeapSys = types.Gauge(ms.HeapSys)
	rm.LastGC = types.Gauge(ms.LastGC)
	rm.Lookups = types.Gauge(ms.Lookups)
	rm.MCacheInuse = types.Gauge(ms.MCacheInuse)
	rm.MCacheSys = types.Gauge(ms.MCacheSys)
	rm.MSpanInuse = types.Gauge(ms.MSpanInuse)
	rm.MSpanSys = types.Gauge(ms.MSpanSys)
	rm.Mallocs = types.Gauge(ms.Mallocs)
	rm.NextGC = types.Gauge(ms.NextGC)
	rm.NumForcedGC = types.Gauge(ms.NumForcedGC)
	rm.NumGC = types.Gauge(ms.NumGC)
	rm.OtherSys = types.Gauge(ms.OtherSys)
	rm.PauseTotalNs = types.Gauge(ms.PauseTotalNs)
	rm.StackInuse = types.Gauge(ms.StackInuse)
	rm.StackSys = types.Gauge(ms.StackSys)
	rm.Sys = types.Gauge(ms.Sys)
	rm.TotalAlloc = types.Gauge(ms.TotalAlloc)
	rm.PollCount += 1
	rand.Seed(time.Now().UnixNano())
	rm.RandomValue = types.Gauge(rand.Float64())
}

var (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	rm := types.RunMetrics{}
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	logrus.SetReportCaller(true)
	logrus.Info(m.Alloc)
	fill(m, &rm)
	logrus.Infof("%+v", rm)
	// return
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChanel
			switch s {
			// kill -SIGHUP XXXX [XXXX - идентификатор процесса для программы]
			case syscall.SIGINT:
				fmt.Println("Signal interrupt triggered.")
				exitChan <- 0
				// kill -SIGTERM XXXX [XXXX - идентификатор процесса для программы]
			case syscall.SIGTERM:
				fmt.Println("Signal terminte triggered.")
				exitChan <- 0

				// kill -SIGQUIT XXXX [XXXX - идентификатор процесса для программы]
			case syscall.SIGQUIT:
				fmt.Println("Signal quit triggered.")
				exitChan <- 0

			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	tickerFill := time.NewTicker(pollInterval)
	go func() {
		for {
			<-tickerFill.C
			runtime.ReadMemStats(&m)
			fill(m, &rm)
		}
	}()

	tickerSendMetrics := time.NewTicker(reportInterval)
	go func() {
		for {
			<-tickerSendMetrics.C
			sendMetrics(rm)
		}
	}()

	// }
	// runtime.ReadMemStats()
	exitCode := <-exitChan
	//stoping ticker
	logrus.Warn("Stopping tickerFill")
	tickerFill.Stop()
	logrus.Warn("Stopping tickerSendMetrics")
	tickerSendMetrics.Stop()
	logrus.Warn("Exiting with code ", exitCode)
	os.Exit(exitCode)

}

func sendMetrics(rm types.RunMetrics) {
	// в формате: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
	// адрес сервиса (как его писать, расскажем в следующем уроке)
	rmMap := structs.Map(rm)
	val := reflect.ValueOf(rm)
	for i := 0; i < val.NumField(); i++ {
		endpoint := "http://127.0.0.1:8080/update"
		// fmt.Println(rm)
		// logrus.Info(val.Type().Field(i).Type.Name())
		// logrus.Info(val.Type().Field(i).Name)
		// logrus.Info(rmMap[val.Type().Field(i).Name])
		// logrus.Info(fmt.Sprintf("%v", rmMap[val.Type().Field(i).Name]))
		// logrus.Info(reflect.TypeOf(rmMap[val.Type().Field(i).Name]))
		endpoint = fmt.Sprintf("%s/%s/%s/%v", endpoint, strings.ToLower(val.Type().Field(i).Type.Name()), val.Type().Field(i).Name, rmMap[val.Type().Field(i).Name])
		// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
		logrus.Info()
		client := &http.Client{}

		request, err := http.NewRequest(http.MethodPost, endpoint, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		request.Header.Add("Content-Type", "text/plain")
		// request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// печатаем код ответа
		fmt.Println("Статус-код ", response.Status)
		defer response.Body.Close()
		// читаем поток из тела ответа

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// и печатаем его
		fmt.Println(string(body))
	}
	logrus.Info(rm.Alloc, rm.PollCount, rm.RandomValue)
}
