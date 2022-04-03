package main

import (
	"log"
	"net/http"

	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/kmx0/devops/cmd/server/storage"
)

func main() {
	store := storage.NewInMemory()
	handlers.SetRepository(store)

	// var handler3 ApiHandler
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", handlers.HandleGauge)
	mux.HandleFunc("/update/counter/", handlers.HandleCounter)
	mux.HandleFunc("/update/unknown/", handlers.)

	addr := "127.0.0.1:8080"
	server := &http.Server{Handler: mux, Addr: addr}
	// http.Handle("/api/", apiHandler)
	// http.Handle("/api/auth", apiAuthHandler)
	// http.HandleFunc("/update/", HandleMetricsFunc)

	log.Fatal(server.ListenAndServe())
}
