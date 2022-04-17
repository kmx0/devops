package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

var (
	rm *types.RunMetrics = &types.RunMetrics{MapMetrics: make(map[string]interface{})}
)

func ReplaceUnused(cfg *config.Config) {
	address := flag.String("a", "127.0.0.1:8080", "Address on server for Sending Metrics ")
	reportInterval := flag.Duration("r", 10000000000, "REPORT_INTERVAL")
	pollInterval := flag.Duration("p", 5000000000, "POLL_INTERVAL")
	flag.Parse()
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		cfg.Address = *address
	}

	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok {
		cfg.ReportInterval = *reportInterval
	}
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok {
		cfg.PollInterval = *pollInterval
	}
}
func main() {
	// rm := types.RunMetrics{}
	logrus.SetReportCaller(true)
	cfg := config.LoadConfig()
	ReplaceUnused(&cfg)
	logrus.Infof("CFG for AGENT %+v", cfg)
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	// logrus.Info(m.Alloc)
	rm.Set(m)
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
				logrus.Info("Signal interrupt triggered.")
				exitChan <- 0
				// kill -SIGTERM XXXX [XXXX - идентификатор процесса для программы]
			case syscall.SIGTERM:
				logrus.Info("Signal terminte triggered.")
				exitChan <- 0

				// kill -SIGQUIT XXXX [XXXX - идентификатор процесса для программы]
			case syscall.SIGQUIT:
				logrus.Info("Signal quit triggered.")
				exitChan <- 0

			default:
				logrus.Info("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	tickerFill := time.NewTicker(cfg.PollInterval)
	go func() {
		for {
			<-tickerFill.C
			runtime.ReadMemStats(&m)
			rm.Set(m)
		}
	}()

	tickerSendMetrics := time.NewTicker(cfg.ReportInterval)
	go func() {
		for {
			<-tickerSendMetrics.C
			// rm.ULock()
			now := time.Now()
			sendMetricsJSON(cfg.Address)
			fmt.Println(time.Since(now))
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

func sendMetricsJSON(address string) {
	// в формате: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
	// адрес сервиса (как его писать, расскажем в следующем уроке)
	// rm.Lock()

	metricsForBody := rm.GetMetrics()
	endpoint := fmt.Sprintf("http://%s/update/", address)
	// logrus.Info(endpoint, metricsForBody)
	// return
	client := &http.Client{}
	for i := 0; i < len(metricsForBody); i++ {

		bodyBytes, err := json.Marshal(metricsForBody[i])
		if err != nil {
			logrus.Error(err)
			continue
		}
		bodyIOReader := bytes.NewReader(bodyBytes)
		request, err := http.NewRequest(http.MethodPost, endpoint, bodyIOReader)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		request.Header.Add("Content-Type", "application/json")
		// request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		response, err := client.Do(request)
		if err != nil {
			logrus.Error("Error on requesting")
			logrus.Error(err)
			// os.Exit(1)
			continue
		}
		// печатаем код ответа
		logrus.Info("Статус-код ", response.Status)
		defer response.Body.Close()
		// читаем поток из тела ответа

		_, err = io.ReadAll(response.Body)
		if err != nil {
			logrus.Error("Error on Reading body")
			logrus.Error(err)
			// os.Exit(1)
			continue
		}
		// и печатаем его
		// fmt.Println(string(body))
	}
}

// func sendMetrics() {
// 	// в формате: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
// 	// адрес сервиса (как его писать, расскажем в следующем уроке)
// 	// rm.Lock()

// 	ret := rm.Get()
// 	// logrus.Info(ret)
// 	// return
// 	for i := 0; i < len(ret); i++ {

// 		client := &http.Client{}

// 		request, err := http.NewRequest(http.MethodPost, ret[i], nil)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		request.Header.Add("Content-Type", "text/plain")
// 		// request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
// 		response, err := client.Do(request)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		// печатаем код ответа
// 		fmt.Println("Статус-код ", response.Status)
// 		defer response.Body.Close()
// 		// читаем поток из тела ответа

// 		body, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		// и печатаем его
// 		fmt.Println(string(body))
// 	}
// }
