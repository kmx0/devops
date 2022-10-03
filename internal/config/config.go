package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Config struct {
	// Адрес, на котором будет запущен сервер
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	// Интервал передачи метрик серверу
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	// Интервал заполнения метрик
	PollInterval time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	// Частота сохранения метрик на стороне сервера
	// Если равно 0, то необходимо сразу записывать в хранилище
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	// Файл для сохранения метрик
	StoreFile string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	// Необходимость восстанавливать данные из хранилища при запуске
	Restore bool `env:"RESTORE" envDefault:"true"`
	// Фактически соль для хэшироварния данных при передаче
	// и проверке при приеме
	Key string `env:"KEY" `
	// параметры базы данных
	// например:
	// "postgres://postgres:postgres@localhost:5432/metrics"
	DBDSN     string `env:"DATABASE_DSN"`
	// Файл с ключами:
	// Публичный для агента
	// И Приватный для сервера
	CryptoKey string `env:"CRYPTO_KEY"`
}

// Парсинг значений из environment или опций запуска.
func LoadConfig() Config {
	logrus.SetReportCaller(true)
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logrus.Error(err)
	}
	return cfg
}

// Установка значений по умолчанию для опций, не укзанных при старте
// для Агента.
func ReplaceUnusedInAgent(cfg *Config) {
	address := flag.String("a", "127.0.0.1:8080", "Address on server for Sending Metrics ")
	reportInterval := flag.Duration("r", 10000000000, "REPORT_INTERVAL")
	pollInterval := flag.Duration("p", 5000000000, "POLL_INTERVAL")
	cryptoKey := flag.String("crypto-key", "", "CRYPTO_KEY")
	key := flag.String("k", "", "KEY for hash")
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
	if _, ok := os.LookupEnv("KEY"); !ok {
		cfg.Key = *key
	}
	if _, ok := os.LookupEnv("CRYPTO_KEY"); !ok {
		cfg.CryptoKey = *cryptoKey
	}
}

// Установка значений по умолчанию для опций, не укзанных при старте
// для Сервера.
func ReplaceUnusedInServer(cfg *Config) {
	//    = flag.Flag("aa", "Address on Listen").Short('a').Default("127.0.0.1:8080").String()
	address := flag.String("a", "127.0.0.1:8080", "Address on Listen")
	restore := flag.Bool("r", true, "restore from file or not")
	storeInterval := flag.Duration("i", 300000000000, "STORE_INTERVAL")
	storeFile := flag.String("f", "/tmp/devops-metrics-db.json", "STORE_FILE")
	dbDSN := flag.String("d", "", "database URI")
	cryptoKey := flag.String("ck", "", "crypto key for cipher")
	key := flag.String("k", "", "KEY for hash")

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
	if _, ok := os.LookupEnv("KEY"); !ok {
		cfg.Key = *key
	}
	if _, ok := os.LookupEnv("CRYPTO_KEY"); !ok {
		cfg.CryptoKey = *cryptoKey
	}
	logrus.Info(cfg.DBDSN)
	logrus.Info(*dbDSN)
	if _, ok := os.LookupEnv("DATABASE_DSN"); !ok {
		// if !strings.Contains(cfg.DBDSN, "incorr") {
		cfg.DBDSN = *dbDSN
	}
	// }
}
