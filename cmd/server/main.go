package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kmx0/devops/cmd/server/handlers"
	"github.com/kmx0/devops/cmd/server/storage"
)

func main() {
	store := storage.NewInMemory()
	handlers.SetRepository(store)
	// logrus.Info(store)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", Handl)
	// r.Get("/", handlers.HandleGauge)

	// r.Post("/update/gauge/", handlers.HandleGauge)
	// r.Post("/update/counter/", handlers.HandleCounter)
	// // r.Get("/update/unknown/", handlers.HandleUnknown)

	addr := "127.0.0.1:8080"

	r.Route("/{update}", func(r chi.Router) {
		r.Post("/", handlers.HandleSl)
		r.Post("/{counter}", handlers.HandleEmptyCounter)
		r.Post("/{counter}/{metric}/{value}", handlers.HandleCounter)
		r.Post("/{gauge}", handlers.HandleEmptyGauge)
		r.Post("/{gauge}/{metric}/{value}", handlers.HandleGauge)
	})

	log.Fatal(http.ListenAndServe(addr, r))
}

func Handl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
