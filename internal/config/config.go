package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	DBDSN string `env:"DATABASE_DSN"`
	// Файл с ключами:
	// Публичный для агента
	// И Приватный для сервера
	CryptoKey string `env:"CRYPTO_KEY"`
	// Передача аргументов через конфиг JSON
	// Имеет наименьший приоритет
	ConfigJSON string `env:"CONFIG"`
	// TRUSTED_SUBNET - доверенная сеть, с которой стоит принимать запросы
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
	GRPC          bool
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
	configJSON := flag.String("c", "", "path to JSON config")
	grpc := flag.Bool("g", false, "enable GRPC transport")

	flag.Parse()

	cfg.GRPC = *grpc

	var cfgJSON ConfigJSON
	if _, ok := os.LookupEnv("CONFIG"); !ok {
		cfg.ConfigJSON = *configJSON
	}

	if _, ok := os.LookupEnv("CONFIG"); ok || isFlagPassed("c") {
		cfgJSON = LoadConfigJSON(cfg.ConfigJSON)
	}
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		if isFlagPassed("c") && !isFlagPassed("a") {
			cfg.Address = cfgJSON.Address
		} else {
			cfg.Address = *address
		}

	}

	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok {
		if isFlagPassed("c") && !isFlagPassed("r") {
			cfg.ReportInterval = cfgJSON.ReportInterval
		} else {

			cfg.ReportInterval = *reportInterval
		}
	}
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok {
		if isFlagPassed("c") && !isFlagPassed("p") {
			cfg.PollInterval = cfgJSON.PollInterval
		} else {

			cfg.PollInterval = *pollInterval
		}
	}
	if _, ok := os.LookupEnv("KEY"); !ok {
		if isFlagPassed("c") && !isFlagPassed("k") {
			cfg.Key = cfgJSON.Key
		} else {

			cfg.Key = *key
		}
	}
	if _, ok := os.LookupEnv("CRYPTO_KEY"); !ok {
		if isFlagPassed("c") && !isFlagPassed("crypto-key") {
			cfg.CryptoKey = cfgJSON.CryptoKey
		} else {
			cfg.CryptoKey = *cryptoKey
		}
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
	configJSON := flag.String("c", "", "path to JSON config")
	trustedSubnet := flag.String("t", "", "trusted subnet")
	grpc := flag.Bool("g", false, "enable GRPC transport")
	flag.Parse()
	cfg.GRPC = *grpc
	var cfgJSON ConfigJSON
	if _, ok := os.LookupEnv("CONFIG"); !ok {
		cfg.ConfigJSON = *configJSON
	}

	if _, ok := os.LookupEnv("CONFIG"); ok || isFlagPassed("c") {
		cfgJSON = LoadConfigJSON(cfg.ConfigJSON)
	}
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		if isFlagPassed("c") && !isFlagPassed("a") {
			cfg.Address = cfgJSON.Address
		} else {
			cfg.Address = *address
		}

	}

	if _, ok := os.LookupEnv("RESTORE"); !ok {
		if isFlagPassed("c") && !isFlagPassed("r") {

			cfg.Restore = cfgJSON.Restore
		} else {

			cfg.Restore = *restore
		}

	}
	if _, ok := os.LookupEnv("STORE_INTERVAL"); !ok {
		if isFlagPassed("c") && !isFlagPassed("i") {
			cfg.StoreInterval = cfgJSON.StoreInterval
		} else {

			cfg.StoreInterval = *storeInterval
		}
	}

	if _, ok := os.LookupEnv("STORE_FILE"); !ok {
		if isFlagPassed("c") && !isFlagPassed("f") {
			cfg.StoreFile = cfgJSON.StoreFile
		} else {

			cfg.StoreFile = *storeFile
		}
	}
	if _, ok := os.LookupEnv("KEY"); !ok {
		if isFlagPassed("c") && !isFlagPassed("k") {
			cfg.Key = cfgJSON.Key
		} else {
			cfg.Key = *key
		}
	}
	if _, ok := os.LookupEnv("CRYPTO_KEY"); !ok {
		if isFlagPassed("c") && !isFlagPassed("ck") {
			cfg.CryptoKey = cfgJSON.CryptoKey
		} else {

			cfg.CryptoKey = *cryptoKey
		}
	}

	if _, ok := os.LookupEnv("TRUSTED_SUBNET"); !ok {
		if isFlagPassed("c") && !isFlagPassed("t") {
			cfg.TrustedSubnet = cfgJSON.TrustedSubnet
		} else {

			cfg.TrustedSubnet = *trustedSubnet
		}
	}

	logrus.Info(cfg.DBDSN)
	logrus.Info(*dbDSN)
	if _, ok := os.LookupEnv("DATABASE_DSN"); !ok {
		// if !strings.Contains(cfg.DBDSN, "incorr") {
		if isFlagPassed("c") && !isFlagPassed("d") {
			cfg.DBDSN = cfgJSON.DBDSN
		} else {

			cfg.DBDSN = *dbDSN
		}
	}
	// }
}

type ConfigJSON struct {
	// Адрес, на котором будет запущен сервер
	Address string `json:"address"`
	// Интервал передачи метрик серверу
	ReportInterval time.Duration `json:"report_interval"`
	// Интервал заполнения метрик
	PollInterval time.Duration `json:"poll_interval"`
	// Частота сохранения метрик на стороне сервера
	// Если равно 0, то необходимо сразу записывать в хранилище
	StoreInterval time.Duration `json:"store_interval"`
	// Файл для сохранения метрик
	StoreFile string `json:"store_file,omitempty"`
	// Необходимость восстанавливать данные из хранилища при запуске
	Restore bool `json:"restore"`
	// Фактически соль для хэшироварния данных при передаче
	// и проверке при приеме
	Key string `json:"key,omitempty" `
	// параметры базы данных
	// например:
	// "postgres://postgres:postgres@localhost:5432/metrics"
	DBDSN string `json:"database_dsn"`
	// Файл с ключами:
	// Публичный для агента
	// И Приватный для сервера
	CryptoKey string `json:"crypto_key"`
	// TRUSTED_SUBNET - доверенная сеть, с которой стоит принимать запросы
	TrustedSubnet string `json:"trusted_subnet"`
	GRPC          bool   `json:"grpc"`
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Парсинг значений из environment или опций запуска.
func LoadConfigJSON(configPath string) ConfigJSON {
	var cfgJSON ConfigJSON
	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &cfgJSON)

	return cfgJSON
}
