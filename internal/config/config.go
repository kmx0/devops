package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval int `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int `env:"POLL_INTERVAL" envDefault:"2"`
}

func LoadConfig() Config {
	logrus.SetReportCaller(true)
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Error(err)
	}
	return cfg
}
