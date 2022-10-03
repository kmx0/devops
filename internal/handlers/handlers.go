package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/crypto"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/storage"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
)

var store repositories.Repository
var cfg config.Config

func SetRepository(s repositories.Repository) {
	store = s
}

func SetupRouter(cf config.Config) (*gin.Engine, *storage.InMemory) {
	store := storage.NewInMemory(cfg)
	cfg = cf
	SetRepository(store)
	var privateKey *rsa.PrivateKey
	var err error
	if cf.CryptoKey != "" {
		privateKey, err = crypto.ReadPrivateKey(cfg.CryptoKey)
		if err != nil {
			logrus.Error(err)
			return nil, nil
		}

	}
	r := gin.New()
	r.Use(gin.Recovery(),
		Compress(),
		Decompress(),
		Decrypt(privateKey),
		gin.Logger())
	//Added for profiling
	//Listen in address:port/debug/pprof
	pprof.Register(r)
	r.POST("/update/gauge/", HandleWithoutID)
	r.POST("/update/counter/", HandleWithoutID)
	r.POST("/update/:typem/:metric/:value", HandleUpdate)
	r.POST("/update/", HandleUpdateJSON)
	r.POST("/updates/", HandleUpdateBatchJSON)
	r.POST("/value/", HandleValueJSON)

	r.GET("/", HandleAllValues)
	r.GET("/ping", HandlePing)
	r.GET("/value/:typem/:metric", HandleValue)
	return r, store
}

// HandleAllValues - Handling Get request by route /
func HandleAllValues(c *gin.Context) {
	mapGauge, mapCounter, _ := store.GetCurrentMetrics()

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, "%+v\n%+v", mapGauge, mapCounter)
}

// HandlePing - Checking DB connection
// Get Request by route /ping
func HandlePing(c *gin.Context) {

	ok := storage.PingDB(c, cfg.DBDSN)
	if ok {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)
		return
	}
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
}

// HandleValue - Get Request metric value by route /value/:mtype/:metric
func HandleValue(c *gin.Context) {
	typeM := c.Param("typem")
	metric := c.Param("metric")
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
}

// HandleValueJSON - Post Request by route /value/
// returning JSON body
func HandleValueJSON(c *gin.Context) {
	logrus.SetReportCaller(true)

	body := c.Request.Body
	defer body.Close()

	decoder := json.NewDecoder(body)
	var metrics types.Metrics

	err := decoder.Decode(&metrics)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	switch metrics.MType {
	case "counter":
		delta, err := store.GetCounterJSON(metrics.ID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		metrics.Delta = &delta
		if cfg.Key != "" {
			metrics.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", metrics.ID, *metrics.Delta), cfg.Key)
		}
		c.JSON(http.StatusOK, metrics)
		return
	case "gauge":
		value, err := store.GetGaugeJSON(metrics.ID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		metrics.Value = &value
		if cfg.Key != "" {

			metrics.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", metrics.ID, *metrics.Value), cfg.Key)
		}
		c.JSON(http.StatusOK, metrics)
		return
	default:
		c.Status(http.StatusBadRequest)

		return

	}
}

func HandleWithoutID(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

// HandleUpdate - Post Request by route /update/:typem/:metric/:value
// Its Update Metric value in Memory and save to Disk
func HandleUpdate(c *gin.Context) {
	logrus.SetReportCaller(true)
	typeM := c.Param("typem")
	metric := c.Param("metric")
	value := c.Param("value")
	if typeM == "counter" || typeM == "gauge" {
		err := store.Update(typeM, metric, value)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), `strconv.ParseInt: parsing`) || strings.Contains(err.Error(), `strconv.ParseFloat: parsing`):
				c.Status(http.StatusBadRequest)
			default:
				c.Status(http.StatusInternalServerError)
			}

		}
		if cfg.StoreInterval == 0 || cfg.DBDSN != "" {

			store.SaveToDisk(cfg)
		}
	} else {
		c.Status(http.StatusNotImplemented)
	}

}

// HandleUpdate - Post Request by route /update/
// In body JSON Metrics struct
// Its Update Metric value in Memory and save to Disk
func HandleUpdateJSON(c *gin.Context) {
	logrus.SetReportCaller(true)
	body := c.Request.Body
	decoder := json.NewDecoder(body)
	var metrics types.Metrics
	err := decoder.Decode(&metrics)
	c.Header("Content-Length", "0")
	if err != nil {
		logrus.Error(err)
		c.Status(http.StatusInternalServerError)
	}
	defer body.Close()
	if metrics.MType == "counter" || metrics.MType == "gauge" {
		err := store.UpdateJSON(cfg.Key, metrics)

		if err != nil {
			logrus.Error(err)

			switch {
			case strings.Contains(err.Error(), `received nil pointer on Delta`) || strings.Contains(err.Error(), `received nil pointer on Value`):
				c.Status(http.StatusBadRequest)
			case strings.Contains(err.Error(), `hash sum not matched`):
				c.Status(http.StatusBadRequest)
			default:
				c.Status(http.StatusInternalServerError)
			}
		} else if cfg.StoreInterval == 0 || cfg.DBDSN != "" {
			store.SaveToDisk(cfg)

		}

	} else {
		c.Status(http.StatusNotImplemented)
	}

}

// HandleUpdateBatchJSON - like HandleUpdateJSON, but come array metrics
func HandleUpdateBatchJSON(c *gin.Context) {
	logrus.SetReportCaller(true)
	body := c.Request.Body
	decoder := json.NewDecoder(body)
	var metrics []types.Metrics
	err := decoder.Decode(&metrics)
	if err != nil {
		logrus.Error(err)
		c.Status(http.StatusInternalServerError)
	}
	for _, v := range metrics {
		defer body.Close()
		if v.MType == "counter" || v.MType == "gauge" {
			err := store.UpdateJSON(cfg.Key, v)
			if err != nil {
				logrus.Error(err)
				switch {
				case strings.Contains(err.Error(), `received nil pointer on Delta`) || strings.Contains(err.Error(), `received nil pointer on Value`) || strings.Contains(err.Error(), `hash sum not matched`):
					c.Status(http.StatusBadRequest)
				default:
					c.Status(http.StatusInternalServerError)
				}

			}
			if cfg.StoreInterval == 0 || cfg.DBDSN != "" {
				store.SaveToDisk(cfg)
			}
		} else {
			c.Status(http.StatusNotImplemented)

		}
	}

}
