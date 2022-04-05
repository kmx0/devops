package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/sirupsen/logrus"
)

var store repositories.Repository

func SetRepository(s repositories.Repository) {
	store = s
}

func HandleEmptyCounter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
func HandleEmptyGauge(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
func HandleCounter(w http.ResponseWriter, r *http.Request) {
	update := chi.URLParam(r, "update")
	if update != "update" {
		w.WriteHeader(http.StatusNotFound)
	}
	counter := chi.URLParam(r, "counter")
	if counter != "counter" {
		w.WriteHeader(http.StatusBadRequest)
	}
	metric := chi.URLParam(r, "metric")

	value := chi.URLParam(r, "value")
	logrus.Info(value)

	err := store.Update("Counter", metric, value)
	if err != nil {
		switch err.Error() {
		case "not such metric name":
			w.WriteHeader(http.StatusBadRequest)
		case `strconv.ParseInt: parsing "none": invalid syntax`:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		// logrus.Error(err)
	}
	w.Header().Set("Content-Type", "text/plain")
	logrus.Info()
	w.WriteHeader(http.StatusOK)
}

func HandleGauge(w http.ResponseWriter, r *http.Request) {
	update := chi.URLParam(r, "update")
	if update != "update" {
		w.WriteHeader(http.StatusNotFound)
	}
	counter := chi.URLParam(r, "gauge")
	if counter != "gauge" {
		w.WriteHeader(http.StatusBadRequest)
	}
	metric := chi.URLParam(r, "metric")

	value := chi.URLParam(r, "value")
	logrus.Info(value)

	err := store.Update("Gauge", metric, value)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not such metric name"):
			w.WriteHeader(http.StatusBadRequest)
		case strings.Contains(err.Error(), `strconv.ParseFloat: parsing`):
			w.WriteHeader(http.StatusBadRequest)
		default:
			// logrus.Info(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	// logrus.Info(storage.MapGauge)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

func HandleUnknown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func HandleSl(w http.ResponseWriter, r *http.Request) {
	logrus.Info("SL")
	w.WriteHeader(http.StatusNotImplemented)
}
