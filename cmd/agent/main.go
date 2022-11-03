package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/crypto"
	metrics_server "github.com/kmx0/devops/internal/metrics_server/client"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	rm           *types.RunMetrics = &types.RunMetrics{MapMetrics: make(map[string]interface{})}
	buildVersion string            = "N/A"
	buildDate    string            = "N/A"
	buildCommit  string            = "N/A"
	publicKey    *rsa.PublicKey
	err          error
	gconn        *grpc.ClientConn
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
	logrus.WithFields(logrus.Fields{
		"build_version": buildVersion,
		"build_date":    buildDate,
		"build_commit":  buildCommit,
	})
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

	exitChan := make(chan struct{})
	go func() {
		for {
			switch <-signalChanel {
			default:
				logrus.Info("Exit Signal")
				exitChan <- struct{}{}
			}
		}
	}()
	if cfg.GRPC {
		gconn, err = grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logrus.Error(err)
			return
		}
	}
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
			if !cfg.GRPC {
				SendMetricsJSON(cfg)
			} else {
				metrics_server.SendMetricsBatch(context.TODO(), gconn, rm.GetMetrics())
			}
		}
	}()

	exitCode := <-exitChan
	// stoping ticker
	logrus.Warn("Stopping tickerFill")
	tickerFill.Stop()
	logrus.Warn("Stopping tickerSendMetrics")
	tickerSendMetrics.Stop()
	logrus.Warn("Sending unsaved metrics")
	if !cfg.GRPC {
		SendMetricsJSON(cfg)
	} else {
		metrics_server.SendMetricsBatch(context.TODO(), gconn, rm.GetMetrics())
	}
	logrus.Warn("Exiting with code ", exitCode)
	os.Exit(0)

}

// sending Metrics use JSON format
func SendMetricsJSON(cfg config.Config) error {
	metricsForBody := rm.GetMetrics()
	endpoint := fmt.Sprintf("http://%s/update/", cfg.Address)
	client := &http.Client{}
	localIP := GetLocalIP()
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
		request.Header.Add("X-Real-IP", localIP)
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

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
