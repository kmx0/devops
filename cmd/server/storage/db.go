package storage

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

func PingDB(urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
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

	fmt.Println("Successfully connected!")
	return true
}
