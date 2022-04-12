package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func main() {

	var cfg Config
	// допишите код здесь
	err := env.Parse(&cfg)
	if err != nil {
		cfg.Address = "127.0.0.1:8080"
		logrus.Infof("Using default values for cfg: %+v", cfg)
		logrus.Error(err)
	}
	r := handlers.SetupRouter()

	log.Fatal(http.ListenAndServe(cfg.Address, r))
}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
