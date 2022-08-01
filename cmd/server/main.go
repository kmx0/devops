package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
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
	cfg := config.LoadConfig()
	config.ReplaceUnusedInServer(&cfg)
	logrus.Infof("CFG for SERVER  %+v", cfg)
	r, sm := handlers.SetupRouter(cfg)
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
	go http.ListenAndServe(cfg.Address, r)
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
