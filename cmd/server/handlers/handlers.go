package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/devops/cmd/server/storage"
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

var store repositories.Repository
var cfg config.Config

func SetRepository(s repositories.Repository) {
	store = s
}

func SetupRouter(cf config.Config) (*gin.Engine, *storage.InMemory) {
	store, err := storage.NewInMemory(cfg)
	if err != nil {
		logrus.Error(err)
	}
	cfg = cf
	SetRepository(store)

	r := gin.New()
	r.Use(gin.Recovery(),
		Compress(),
		Decompress(),
		gin.Logger())

	r.POST("/update/gauge/", HandleWithoutID)
	r.POST("/update/counter/", HandleWithoutID)
	r.POST("/update/:typem/:metric/:value", HandleUpdate)
	r.POST("/update/", HandleUpdateJSON)
	// r.POST("/update/", HandleUpdateJSON)
	r.POST("/value/", HandleValueJSON)

	r.GET("/", HandleAllValues)
	r.GET("/value/:typem/:metric", HandleValue)
	// r.GET("/value/counter/:metric", HandleValue)
	// r.GET("/", HandleValue)
	// http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ> (со статусом http.StatusOK).
	return r, store
}

func HandleAllValues(c *gin.Context) {
	gg, cntr, _ := store.GetCurrentMetrics()

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, "%+v\n%+v", gg, cntr)
	// ntf("%s/%s/%s/%v", endpoint, val.Type().Field(i).Type.Nam
}

func HandleValue(c *gin.Context) {
	typeM := c.Param("typem")
	metric := c.Param("metric")
	logrus.Info(metric)
	switch typeM {
	case "counter":
		value, err := store.GetCounter(typeM, metric)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, value.String())
		return
	case "gauge":
		value, err := store.GetGauge(typeM, metric)
		logrus.Info(err)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, value.String())
		return
	default:
		c.Status(http.StatusNotFound)

		return

	}
	// c.Status(http.StatusNotFound)
	// c.String(http)
}

func HandleValueJSON(c *gin.Context) {
	logrus.SetReportCaller(true)
	// logrus.Info(typeM, metric, value)

	// var bodyBytes []byte
	body := c.Request.Body
	// a := []byte{}
	// _, err := body.Read(a)
	// logrus.Infof("%+v", string(a))
	// if err != nil {
	// 	c.Status(http.StatusBadRequest)
	// }
	defer body.Close()

	decoder := json.NewDecoder(body)
	// var t test_struct
	var metrics types.Metrics

	err := decoder.Decode(&metrics)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	// c.Request.Bodyjj
	logrus.Info(metrics)
	switch metrics.MType {
	case "counter":
		value, err := store.GetCounterJSON(metrics)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, value)
		// c.String(http.StatusOK, value.String())
		return
	case "gauge":
		value, err := store.GetGaugeJSON(metrics)
		logrus.Info(err)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, value)
		return
	default:
		c.Status(http.StatusBadRequest)

		return

	}
	// c.Status(http.StatusNotFound)
	// c.String(http)
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
		if cfg.StoreInterval == 0 {

			store.SaveToDisk(cfg)
		}
	case "gauge":

		err := store.Update(typeM, metric, value)
		logrus.Info(err)

		if err != nil {
			switch {
			// case strings.Contains(err.Error(), `not such metric`):
			// 	c.Status(http.StatusBadRequest)
			case strings.Contains(err.Error(), `strconv.ParseFloat: parsing`):
				c.Status(http.StatusBadRequest)
			default:
				// logrus.Info(err)
				c.Status(http.StatusInternalServerError)
			}
		}
		// c.Status(http.StatusOK)
		if cfg.StoreInterval == 0 {

			store.SaveToDisk(cfg)
		}
	default:
		c.Status(http.StatusNotImplemented)
	}

}

func HandleUpdateJSON(c *gin.Context) {
	logrus.SetReportCaller(true)
	// logrus.Info(typeM, metric, value)

	// var bodyBytes []byte
	// logrus.Info("UPDATEJSON!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1")
	body := c.Request.Body

	decoder := json.NewDecoder(body)
	// var t test_struct
	var metrics types.Metrics
	err := decoder.Decode(&metrics)
	if err != nil {
		logrus.Error(err)
		// panic(err)
		c.Status(http.StatusInternalServerError)
	}
	defer body.Close()
	// logrus.Info(metrics)
	// logrus.Info("UPDATE")
	switch metrics.MType {
	case "counter":
		err := store.UpdateJSON(metrics)
		if err != nil {
			logrus.Info(err)
			switch {
			case strings.Contains(err.Error(), `recieved nil pointer on Delta`):
				c.Status(http.StatusBadRequest)
			default:
				c.Status(http.StatusInternalServerError)
			}

			// logrus.Error(err)
		}
		// c.Status(http.StatusOK)
		if cfg.StoreInterval == 0 {
			store.SaveToDisk(cfg)
		}

	case "gauge":

		err := store.UpdateJSON(metrics)
		// logrus.Info(err)
		if err != nil {
			switch {
			// case strings.Contains(err.Error(), `not such metric`):
			// 	c.Status(http.StatusBadRequest)
			case strings.Contains(err.Error(), `recieved nil pointer on Value`):
				c.Status(http.StatusBadRequest)
			default:
				// logrus.Info(err)
				c.Status(http.StatusInternalServerError)
			}
		}
		// c.Status(http.StatusOK)
		if cfg.StoreInterval == 0 {
			store.SaveToDisk(cfg)
		}
	default:
		c.Status(http.StatusNotImplemented)
	}

}
