package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func PingDB(urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
	logrus.Info(urlExample)
	db, err := sql.Open("postgres", urlExample)
	if err != nil {
		logrus.Error(err)
		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logrus.Error(err)
		return false
	}

	logrus.Info("Successfully connected!")
	return true
}
