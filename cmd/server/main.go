package main

import (
	"log"
	"net/http"

	"github.com/kmx0/devops/cmd/server/handlers"
)

func main() {

	r := handlers.NewRouter()
	addr := "127.0.0.1:8080"

	log.Fatal(http.ListenAndServe(addr, r))
}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
