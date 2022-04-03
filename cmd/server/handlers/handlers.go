package handlers

import (
	"net/http"
	"strings"

	"github.com/kmx0/devops/internal/repositories"
	"github.com/sirupsen/logrus"
)

var store repositories.Repository

func SetRepository(s repositories.Repository) {
	store = s
}

func HandleCounter(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	logrus.Info(url)
	fields := strings.Split(url, "/")
	err := store.Update("Counter", fields[3], fields[4])
	if err != nil {
		logrus.Error(err)
		
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func HandleGauge(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	logrus.Info(url)
	fields := strings.Split(url, "/")
	err := store.Update("Gauge", fields[3], fields[4])
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// logrus.Info(storage.MapGauge)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
