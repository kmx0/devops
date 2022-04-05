package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/devops/cmd/server/storage"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/sirupsen/logrus"
)

var store repositories.Repository

func SetRepository(s repositories.Repository) {
	store = s
}

func SetupRouter() *gin.Engine {
	store := storage.NewInMemory()
	SetRepository(store)

	r := gin.Default()
	r.POST("/update/gauge/", HandleWithoutID)
	r.POST("/update/counter/", HandleWithoutID)
	r.POST("/update/:typem/:metric/:value", HandleUpdate)
	return r
}
func HandleWithoutID(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

func HandleUpdate(c *gin.Context) {
	logrus.SetReportCaller(true)
	typeM := c.Param("typem")
	metric := c.Param("metric")
	value := c.Param("value")
	logrus.Info(typeM, metric, value)
	switch typeM {
	case "counter":
		err := store.Update(typeM, metric, value)
		if err != nil {
			logrus.Info(err)
			switch {
			case strings.Contains(err.Error(), `strconv.ParseInt: parsing`):
				c.Status(http.StatusBadRequest)
			default:
				c.Status(http.StatusInternalServerError)
			}

			// logrus.Error(err)
		}
		// c.Status(http.StatusOK)
	case "gauge":
		logrus.Info("gauge")

		err := store.Update(typeM, metric, value)
		logrus.Info(err)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), `strconv.ParseFloat: parsing`):
				c.Status(http.StatusBadRequest)
			default:
				// logrus.Info(err)
				c.Status(http.StatusInternalServerError)
			}
		}
		// c.Status(http.StatusOK)
	default:
		c.Status(http.StatusNotImplemented)
	}

}

// func HandleUnknown(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusNotImplemented)
// }

// func HandleSl(w http.ResponseWriter, r *http.Request) {
// 	logrus.Info("SL")
// 	w.WriteHeader(http.StatusNotImplemented)
// }
