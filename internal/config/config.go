package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
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

func ReplaceUnusedInAgent(cfg *Config) {
	address := flag.String("a", "127.0.0.1:8080", "Address on server for Sending Metrics ")
	reportInterval := flag.Duration("r", 10000000000, "REPORT_INTERVAL")
	pollInterval := flag.Duration("p", 5000000000, "POLL_INTERVAL")
	flag.Parse()
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		cfg.Address = *address
	}

	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok {
		cfg.ReportInterval = *reportInterval
	}
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok {
		cfg.PollInterval = *pollInterval
	}
}

func ReplaceUnusedInServer(cfg *Config) {
	//    = flag.Flag("aa", "Address on Listen").Short('a').Default("127.0.0.1:8080").String()
	address := flag.String("a", "127.0.0.1:8080", "Address on Listen")
	restore := flag.Bool("r", true, "restore from file or not")
	storeInterval := flag.Duration("i", 300000000000, "STORE_INTERVAL")
	storeFile := flag.String("f", "/tmp/devops-metrics-db.json", "STORE_FILE")

	flag.Parse()

	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		cfg.Address = *address
	}
	if _, ok := os.LookupEnv("RESTORE"); !ok {

		cfg.Restore = *restore
	}
	if _, ok := os.LookupEnv("STORE_INTERVAL"); !ok {

		cfg.StoreInterval = *storeInterval
	}
	if _, ok := os.LookupEnv("STORE_FILE"); !ok {
		cfg.StoreFile = *storeFile
	}
}
