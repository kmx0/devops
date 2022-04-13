package main

import (
	"log"
	"net/http"

	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/kmx0/devops/internal/config"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func main() {
	logrus.SetReportCaller(true)

	cfg := config.LoadConfig()
	logrus.Infof("CFG for SERVER  %+v", cfg)
	r := handlers.SetupRouter()

	log.Fatal(http.ListenAndServe(cfg.Address, r))
}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
