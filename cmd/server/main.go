package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/kmx0/devops/internal/config"
	"github.com/sirupsen/logrus"
)

func ReplaceUnused(cfg *config.Config) {
	//    = flag.Flag("aa", "Address on Listen").Short('a').Default("127.0.0.1:8080").String()
	address := flag.String("a", "127.0.0.1:8080", "Address on Listen")
	restore := flag.Bool("r", true, "restore from file or not")
	storeInterval := flag.Duration("i", 300000000000, "STORE_INTERVAL")
	storeFile := flag.String("f", "/tmp/devops-metrics-db.json", "STORE_FILE")

	flag.Parse()

	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		cfg.Address = *address
	}
	if _, ok := os.LookupEnv("RESTORE"); !ok {

		cfg.Restore = *restore
	}
	if _, ok := os.LookupEnv("STORE_INTERVAL"); !ok {

		cfg.StoreInterval = *storeInterval
	}
	if _, ok := os.LookupEnv("STORE_FILE"); !ok {
		cfg.StoreFile = *storeFile
	}
}
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
	ReplaceUnused(&cfg)
	logrus.Infof("CFG for SERVER  %+v", cfg)
	r, sm := handlers.SetupRouter(cfg)
	if cfg.Restore {
		sm.RestoreFromDisk(cfg)
	}
	// logrus.Infof("%+v", sm.ArrayJSONMetrics)
	tickerStore := time.NewTicker(cfg.StoreInterval)
	if cfg.StoreInterval != 0 {

		go func() {
			for {
				<-tickerStore.C
				// runtime.ReadMemStats(&m)
				// rm.Set(m)
				// metrics := []types.Metrics{}

				// metrics := storage.ConvertMapsToMetrisc(sm)
				// sm.WriteMetrics(&metrics)
				// sm.Close()
				sm.SaveToDisk(cfg)
				logrus.Infof("Saving data to file %s", cfg.StoreFile)
			}
		}()
	}
	go http.ListenAndServe(cfg.Address, r)
	logrus.Info("EFDVfdvfvfvfewv!!!!!!!!!!!!!!!!!!!!1")
	exitCode := <-exitChan
	//stoping ticker
	logrus.Warn("Stopping tickerStore")

	tickerStore.Stop()

	logrus.Warn("Saving last data")
	sm.SaveToDisk(cfg)

	logrus.Warn("Exiting with code ", exitCode)
	os.Exit(exitCode)
}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
