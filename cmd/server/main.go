package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/gin-gonic/gin"
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/handlers"
	rpc "github.com/kmx0/devops/internal/rpc/server"
	"github.com/kmx0/devops/internal/storage"
	"github.com/sirupsen/logrus"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
	sm           *storage.InMemory
	r            *gin.Engine
)

func main() {
	logrus.SetReportCaller(true)
	logrus.WithFields(logrus.Fields{
		"build_version": buildVersion,
		"build_date":    buildDate,
		"build_commit":  buildCommit,
	})
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
			default:
				logrus.Info("Exit Signal for server")
				exitChan <- 0
			}
		}
	}()
	cfg := config.LoadConfig()
	config.ReplaceUnusedInServer(&cfg)
	logrus.Infof("CFG for SERVER  %+v", cfg)
	if !cfg.GRPC {
		r, sm = handlers.SetupRouter(cfg)
	} else {
		sm = storage.NewInMemory(cfg)
		rpc.NewRPCServer(cfg, sm, cfg.Address)
	}
	if cfg.Restore {
		sm.RestoreFromDisk(cfg)
	}
	tickerStore := time.NewTicker(cfg.StoreInterval)
	if cfg.StoreInterval != 0 {

		go func() {
			for {
				<-tickerStore.C

				sm.SaveToDisk(cfg)
				logrus.Infof("Saving data to file %s", cfg.StoreFile)

			}
		}()
	}
	if !cfg.GRPC {

		go http.ListenAndServe(cfg.Address, r)
	}
	exitCode := <-exitChan
	//stoping ticker
	logrus.Warn("Stopping tickerStore")

	tickerStore.Stop()

	logrus.Warn("Saving last data")
	sm.SaveToDisk(cfg)
	// globalCtx.Done()
	logrus.Warn("Exiting with code ", exitCode)
	os.Exit(exitCode)
}
