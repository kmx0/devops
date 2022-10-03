package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
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
	rm           *types.RunMetrics = &types.RunMetrics{MapMetrics: make(map[string]interface{})}
	buildVersion string
	buildDate    string
	buildCommit  string
	publicKey    *rsa.PublicKey
	err          error
)

func main() {
	logrus.SetReportCaller(true)
	cfg := config.LoadConfig()
	config.ReplaceUnusedInAgent(&cfg)
	if cfg.CryptoKey != "" {
		publicKey, err = crypto.ReadPublicKey(cfg.CryptoKey)
		if err != nil {
			logrus.Error(err)
			return
		}

	}
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)
	logrus.Infof("CFG for AGENT %+v", cfg)
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	rm.Set(m)
	rm.SetGopsutil()
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			switch <-signalChanel {
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
	// every cfg.PollInterval seconds will filled RuntimeMemory values
	tickerFill := time.NewTicker(cfg.PollInterval)
	go func() {
		for {
			<-tickerFill.C
			runtime.ReadMemStats(&m)
			rm.Set(m)
			rm.SetGopsutil()
		}
	}()
	// ever cfg.ReportInterval seconds will sended to server Metrics
	tickerSendMetrics := time.NewTicker(cfg.ReportInterval)
	go func() {
		for {
			<-tickerSendMetrics.C
			now := time.Now()
			SendMetricsJSON(cfg)
			logrus.Info(time.Since(now))
		}
	}()

	exitCode := <-exitChan
	// stoping ticker
	logrus.Warn("Stopping tickerFill")
	tickerFill.Stop()
	logrus.Warn("Stopping tickerSendMetrics")
	tickerSendMetrics.Stop()
	logrus.Warn("Exiting with code ", exitCode)
	os.Exit(exitCode)

}

// sending Metrics use JSON format
func SendMetricsJSON(cfg config.Config) error {
	metricsForBody := rm.GetMetrics()
	endpoint := fmt.Sprintf("http://%s/update/", cfg.Address)
	client := &http.Client{}

	for i := 0; i < len(metricsForBody); i++ {
		if cfg.Key != "" {
			AddHash(cfg.Key, &metricsForBody[i])
		}
		bodyBytes, err := json.Marshal(metricsForBody[i])
		if err != nil {
			logrus.Error(err)
			continue
		}
		if cfg.CryptoKey != "" {
			bodyBytes, err = crypto.EncryptData(*publicKey, bodyBytes)
			if err != nil {
				logrus.Error(err)
				break
			}
		}
		bodyIOReader := bytes.NewReader(bodyBytes)
		request, err := http.NewRequest(http.MethodPost, endpoint, bodyIOReader)
		errors.Is(nil, err)
		if err != nil {
			logrus.Error(err)
			return err
		}

		request.Header.Add("Content-Type", "application/json")
		// if cfg.CryptoKey != "" {
		// 	request.Header.Add("Encrypted", "true")

		// }
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
			continue
		}
	}
	return nil
}

// add Hash to JSON Metrics
func AddHash(key string, metricsP *types.Metrics) error {

	switch metricsP.MType {
	case "counter":
		metricsP.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", metricsP.ID, *metricsP.Delta), key)
		return nil
	case "gauge":
		metricsP.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", metricsP.ID, *metricsP.Value), key)
		return nil
	}
	return errors.New("unknown metric type")
}
