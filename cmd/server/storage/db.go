package storage

import (
	"context"
	"database/sql"

	"github.com/kmx0/devops/internal/types"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB
var DBName = "metrics"
var TableName = "praktikum"

func PingDB(ctx context.Context, urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
	logrus.Info(urlExample)
	var err error
	DB, err = sql.Open("postgres", urlExample)
	if err != nil {
		logrus.Error(err)
		return false
	}
	// defer DB.Close()

	err = DB.Ping()
	if err != nil {
		logrus.Error(err)
		return false
	}

	logrus.Info("Successfully connected!")
	if !CheckTableExist() {
		AddTabletoDB()
	}
	logrus.Info(CheckTableExist())
	return true
}

func CheckDBExist() bool {
	// DB.Begin()
	// create database test
	listDB := `SELECT datname FROM pg_database;`
	rows, err := DB.QueryContext(context.Background(), listDB)
	if err != nil {
		logrus.Error(err)
	}
	// c, _ := result
	defer rows.Close()
	for rows.Next() {
		var res string
		rows.Scan(&res)
		logrus.Info(res)
		if res == DBName {
			return true
		}
	}
	err = rows.Err()
	if err != nil {
		return false
	}
	return false
	// logrus.Infof("%+v", res)

	// // dynamic
	// insertDynStmt := `insert into "Students"("Name", "Roll") values($1, $2)`
	// _, e = DB.Exec(insertDynStmt, "Jane", 2)
	// CheckError(e)
}

func CheckTableExist() bool {
	// DB.Begin()
	// create database test
	listDB := `SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema='public';`
	rows, err := DB.Query(listDB)
	if err != nil {
		logrus.Error(err)
	}
	// c, _ := result
	defer rows.Close()
	for rows.Next() {
		var res string
		rows.Scan(&res)
		logrus.Info(res)
		if res == TableName {
			return true
		}
	}
	err = rows.Err()
	if err != nil {
		return false
	}
	return false
	// logrus.Infof("%+v", res)

	// // dynamic
	// insertDynStmt := `insert into "Students"("Name", "Roll") values($1, $2)`
	// _, e = DB.Exec(insertDynStmt, "Jane", 2)
	// CheckError(e)
}

func AddTabletoDB() {

	req := `CREATE TABLE praktikum (
		ID varchar(255),
		Type varchar(255),
		Delta numeric,
		Value double precision
	);`
	rows, err := DB.Query(req)
	if err != nil {
		logrus.Error(err)
	}
	// c, _ := result
	defer rows.Close()
	for rows.Next() {
		var res string
		rows.Scan(&res)
		logrus.Info(res)

	}
	err = rows.Err()
	if err != nil {
		logrus.Error(err)
	}

}

func SaveDataToDB(sm *InMemory) {
	sm.Lock()
	defer sm.Unlock()
	// TRUNCATE TABLE COMPANY
	// metrics := make([]types.Metrics, len(sm.MapCounter)+len(sm.MapGauge))

	keysCounter := make([]string, 0, len(sm.MapCounter))
	keysGauge := make([]string, 0, len(sm.MapGauge))

	for k := range sm.MapCounter {
		keysCounter = append(keysCounter, k)

	}

	for k := range sm.MapGauge {
		keysGauge = append(keysGauge, k)

	}
	for i := 0; i < len(keysCounter); i++ {
		insertCounter := `INSERT INTO praktikum(ID, Type, Delta) values($1, $2, $3)`
		_, err := DB.Exec(insertCounter, keysCounter[i], "counter", int(sm.MapCounter[keysCounter[i]]))
		if err != nil {
			logrus.Error(err)
		}
	}
	for i := 0; i < len(keysGauge); i++ {
		insertGauge := `INSERT INTO praktikum(ID, Type, Value) values($1, $2, $3)`
		_, err := DB.Exec(insertGauge, keysGauge[i], "gauge", float64(sm.MapGauge[keysGauge[i]]))
		if err != nil {
			logrus.Error(err)
		}
	}

	// listDB := `SELECT datname FROM pg_database;`

}
func RestoreDataFromDB(sm *InMemory) {
	sm.Lock()
	defer sm.Unlock()
	err := DB.Ping()
	if err != nil {
		logrus.Error(err)
		return
	}
	ctx := context.Background()
	listCounter := "SELECT ID, Delta FROM praktikum WHERE Type='counter';"
	rowsC, err := DB.QueryContext(ctx, listCounter)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer rowsC.Close()
	for rowsC.Next() {
		var id string
		var delta int64
		rowsC.Scan(&id, &delta)
		logrus.Info(id)
		logrus.Info(delta)
		sm.MapCounter[id] = types.Counter(delta)
	}

	err = rowsC.Err()
	if err != nil {
		logrus.Error(err)
	}
	listGauge := `SELECT ID, Value FROM praktikum WHERE Type='gauge';`
	rowsG, err := DB.Query(listGauge)
	if err != nil {
		logrus.Error(err)
	}
	// c, _ := result
	defer rowsG.Close()
	for rowsG.Next() {
		var id string
		var value float64
		rowsG.Scan(&id, &value)
		logrus.Info(id)
		logrus.Info(value)
		sm.MapGauge[id] = types.Gauge(value)
	}
	err = rowsG.Err()
	if err != nil {
		logrus.Error(err)
	}
}
