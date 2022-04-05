package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kmx0/devops/cmd/server/storage"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/sirupsen/logrus"
)

var store repositories.Repository

func SetRepository(s repositories.Repository) {
	store = s
}

// func HandleEmptyCounter(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusBadRequest)
// }
// func HandleEmptyGauge(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusBadRequest)
// }
func NewRouter() chi.Router {
	store := storage.NewInMemory()
	SetRepository(store)
	// logrus.Info(store)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// r.Get("/", Handl)
	// r.Get("/", handlers.HandleGauge)

	// r.Post("/update/gauge/", handlers.HandleGauge)
	// r.Post("/update/counter/", handlers.HandleCounter)
	// // r.Get("/update/unknown/", handlers.HandleUnknown)

	r.Post("/update/{typem}/{metric}/{value}", HandleUpdate)
return r
}
func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	logrus.Info(r.URL.String())
	typeM := chi.URLParam(r, "typem")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")
	logrus.Info(typeM)
	logrus.Info(metric)
	logrus.Info(value)
	switch typeM {
	case "counter":
		logrus.Info("counter")

		err := store.Update(typeM, metric, value)
		if err != nil {
			switch err.Error() {
			case `strconv.ParseInt: parsing "none": invalid syntax`:
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			// logrus.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	case "gauge":
		logrus.Info("gauge")

		err := store.Update(typeM, metric, value)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), `strconv.ParseFloat: parsing`):
				w.WriteHeader(http.StatusBadRequest)
			default:
				// logrus.Info(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		w.WriteHeader(http.StatusOK)
	default:
		logrus.Info("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	}

}

// func HandleUnknown(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusNotImplemented)
// }

// func HandleSl(w http.ResponseWriter, r *http.Request) {
// 	logrus.Info("SL")
// 	w.WriteHeader(http.StatusNotImplemented)
// }
