package main

import (
	"log"
	"net/http"

	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/kmx0/devops/internal/config"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func main() {

	cfg := config.LoadConfig("server")

	r := handlers.SetupRouter()

	log.Fatal(http.ListenAndServe(cfg.Address, r))
}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
