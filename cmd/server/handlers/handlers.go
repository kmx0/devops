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
	logrus.Info(len(fields))
	switch len(fields) {
	case 5:
		err := store.Update("Counter", fields[3], fields[4])
		if err != nil {
			switch err.Error() {
			case "not such metric name":
				w.WriteHeader(http.StatusBadRequest)
			case `strconv.ParseInt: parsing "none": invalid syntax`:
				w.WriteHeader(http.StatusBadRequest)
			default:
				// logrus.Info(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			// logrus.Error(err)
		}
	case 4:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func HandleGauge(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	logrus.Info(url)
	fields := strings.Split(url, "/")
	switch len(fields) {
	case 5:
		err := store.Update("Gauge", fields[3], fields[4])
		if err != nil {
			// logrus.Error(err)
			switch{
			case strings.Contains(err.Error(), "not such metric name"):
				w.WriteHeader(http.StatusBadRequest)
			case strings.Contains(err.Error(), `strconv.ParseFloat: parsing`):
				w.WriteHeader(http.StatusBadRequest)
			default:
				// logrus.Info(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			// logrus.Error(err)
		}
	case 4:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	// logrus.Info(storage.MapGauge)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func HandleUnknown(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
}