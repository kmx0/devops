package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/crypto"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

var (
	rm *types.RunMetrics = &types.RunMetrics{MapMetrics: make(map[string]interface{})}
)

func main() {
	logrus.SetReportCaller(true)
	cfg := config.LoadConfig()
	config.ReplaceUnusedInAgent(&cfg)
	logrus.Infof("CFG for AGENT %+v", cfg)
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	rm.Set(m)
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
			now := time.Now()
			sendMetricsJSON(cfg)
			fmt.Println(time.Since(now))
		}
	}()

	exitCode := <-exitChan
	//stoping ticker
	logrus.Warn("Stopping tickerFill")
	tickerFill.Stop()
	logrus.Warn("Stopping tickerSendMetrics")
	tickerSendMetrics.Stop()
	logrus.Warn("Exiting with code ", exitCode)
	os.Exit(exitCode)

}

// func sendMetricsBatchJSON(cfg config.Config) {
// 	metricsForBody := rm.GetMetrics()
// 	endpoint := fmt.Sprintf("http://%s/updates/", cfg.Address)
// 	client := &http.Client{}

// 	logrus.Info(cfg)
// 	for i := 0; i < len(metricsForBody); i++ {
// 		if cfg.Key != "" {
// 			AddHash(cfg.Key, &metricsForBody[i])
// 		}
// 		// logrus.Infof()
// 	}
// 	bodyBytes, err := json.Marshal(metricsForBody)
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	bodyIOReader := bytes.NewReader(bodyBytes)
// 	request, err := http.NewRequest(http.MethodPost, endpoint, bodyIOReader)
// 	if err != nil {
// 		logrus.Error(err)
// 		os.Exit(1)
// 	}

// 	request.Header.Add("Content-Type", "application/json")
// 	response, err := client.Do(request)
// 	if err != nil {
// 		logrus.Error("Error on requesting")
// 		logrus.Error(err)
// 	}
// 	// печатаем код ответа
// 	logrus.Info("Статус-код ", response.Status)
// 	defer response.Body.Close()
// 	// читаем поток из тела ответа
// 	_, err = io.ReadAll(response.Body)
// 	if err != nil {
// 		logrus.Error("Error on Reading body")
// 		logrus.Error(err)
// 		// os.Exit(1)
// 	}

// }

func sendMetricsJSON(cfg config.Config) {
	metricsForBody := rm.GetMetrics()
	endpoint := fmt.Sprintf("http://%s/update/", cfg.Address)
	client := &http.Client{}

	logrus.Info(cfg)
	for i := 0; i < len(metricsForBody); i++ {
		if cfg.Key != "" {
			AddHash(cfg.Key, &metricsForBody[i])
		}
		// logrus.Infof()
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
		response, err := client.Do(request)
		if err != nil {
			logrus.Error("Error on requesting")
			logrus.Error(err)
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
	}
}

func AddHash(key string, metricsP *types.Metrics) {

	switch metricsP.MType {
	case "counter":
		metricsP.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", metricsP.ID, *metricsP.Delta), key)
	case "gauge":
		metricsP.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", metricsP.ID, *metricsP.Value), key)
	}
}
