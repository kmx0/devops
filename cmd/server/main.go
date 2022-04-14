package main

import (
	"net/http"
	"time"

	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/kmx0/devops/internal/config"
	"github.com/sirupsen/logrus"
)

// type Config struct {
// 	Address string `env:"ADDRESS"`
// }

func main() {
	logrus.SetReportCaller(true)

	cfg := config.LoadConfig()
	logrus.Infof("CFG for SERVER  %+v", cfg)
	r, sm := handlers.SetupRouter(cfg)
	if cfg.Restore {
		sm.RestoreFromDisk(cfg)
	}
	if cfg.StoreInterval != 0 {
		tickerStore := time.NewTicker(cfg.StoreInterval)

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
	logrus.Fatal(http.ListenAndServe(cfg.Address, r))

}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
