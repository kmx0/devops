package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/kmx0/devops/internal/types"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB
var Conn *pgx.Conn
var DBName = "metrics"
var TableName = "praktikum"

func PingDB(ctx context.Context, urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
	logrus.Info(urlExample)
	var err error
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	Conn, err = pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())

	err = Conn.Ping(context.Background())
	if err != nil {
		logrus.Error(err)
		return false
	}

	logrus.Info("Successfully connected!")
	logrus.Info(CheckDBExist())
	if !CheckTableExist() {
		AddTabletoDB()
	}
	// logrus.Info(CheckTableExist())
	return true
}

func CheckDBExist() bool {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return false
	}
	listDB := `SELECT datname FROM pg_database;`
	rows, err := Conn.Query(context.Background(), listDB)
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
}

func CheckTableExist() bool {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return false
	}
	listTables := `SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema='public';`
	rows, err := Conn.Query(context.Background(), listTables)
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
}

func AddTabletoDB() {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return
	}
	req := `CREATE TABLE praktikum (
		ID varchar(255) UNIQUE,
		Type varchar(255),
		Delta numeric,
		Value double precision
	);`
	rows, err := Conn.Query(context.Background(), req)
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
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return
	}
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
		_, err := Conn.Exec(context.Background(), insertCounter, keysCounter[i], "counter", int(sm.MapCounter[keysCounter[i]]))
		if err != nil {
			updateCounter := `UPDATE praktikum SET Type = $1, Delta = $2 WHERE ID = $3;`
			_, err := Conn.Exec(context.Background(), updateCounter, "counter", int(sm.MapCounter[keysCounter[i]]), keysCounter[i])
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	for i := 0; i < len(keysGauge); i++ {
		insertGauge := `INSERT INTO praktikum(ID, Type, Value) values($1, $2, $3)`
		_, err := Conn.Exec(context.Background(), insertGauge, keysGauge[i], "gauge", float64(sm.MapGauge[keysGauge[i]]))
		if err != nil {
			updateGauge := `UPDATE praktikum SET Type = $1, Value = $2 WHERE ID = $3;`
			_, err := Conn.Exec(context.Background(), updateGauge, "gauge", float64(sm.MapGauge[keysGauge[i]]), keysGauge[i])
			if err != nil {
				logrus.Error(err)
			}
		}
	}

}
func RestoreDataFromDB(sm *InMemory) {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return
	}
	sm.Lock()
	defer sm.Unlock()
	// err := Conn.Ping()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }
	ctx := context.Background()
	listCounter := "SELECT ID, Delta FROM praktikum WHERE Type='counter';"
	rowsC, err := Conn.Query(ctx, listCounter)
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
	rowsG, err := Conn.Query(ctx, listGauge)
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
