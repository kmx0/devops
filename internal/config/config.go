package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func LoadServer(cfg Config) Config {
	// logrus.SetReportCaller(true)
	err := env.Parse(&cfg)
	if err != nil {
		cfg.Address = "127.0.0.1:8080"
		logrus.Infof("Using default values for cfg: %+v", cfg)
		logrus.Error(err)
	}
	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:8080"
		logrus.Infof("Using default values for Address: %+v", cfg.Address)
	}
	logrus.Infof("%+v", cfg)
	return cfg
}

func LoadAgent(cfg Config) Config {
	// logrus.SetReportCaller(true)
	err := env.Parse(&cfg)
	if err != nil {
		cfg.Address = "127.0.0.1:8080"
		cfg.ReportInterval = 10
		cfg.PollInterval = 2
		logrus.Infof("Using default values for cfg: %+v", cfg)
		logrus.Error(err)
	}
	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:8080"
		logrus.Infof("Using default values for Address: %+v", cfg.Address)
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 2
		logrus.Infof("Using default values for PollInterval: %+v", cfg.PollInterval)
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = 10
		logrus.Infof("Using default values for ReportInterval: %+v", cfg.ReportInterval)
	}
	return cfg
}
func LoadConfig(mode string) Config {
	var cfg Config
	// допишите код здесь
	switch mode {
	case "agent":
		return LoadAgent(cfg)
	case "server":
		return LoadServer(cfg)
	default:
		return LoadAgent(cfg)
	}

}
