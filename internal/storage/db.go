package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/kmx0/devops/internal/types"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB
var Conn *pgx.Conn
var DBName = "metrics"
var TableName = "praktikum"

// PingDB - Connect to DB, check Exist tables and add new tables
func PingDB(ctx context.Context, urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
	var err error
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	Conn, err = pgx.Connect(context.Background(), urlExample)
	if err != nil {
		logrus.Errorf("unable to connect to database: %v\n", err)
		return false
	}
	err = Conn.Ping(context.Background())
	if err != nil {
		logrus.Error(err)
		return false
	}
	if !CheckTableExist() {
		AddTabletoDB()
	}
	return true
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
	defer rows.Close()
	for rows.Next() {
		var res string
		err := rows.Scan(&res)
		if err != nil {
			return false
		} else if res == TableName {
			return true

		}
	}
	err = rows.Err()
	if err != nil {
		return false
	}
	return false
}

func AddTabletoDB() bool {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return false
	}
	req := `CREATE TABLE praktikum (
		ID varchar(255) UNIQUE,
		Type varchar(255),
		Delta numeric,
		Value double precision
	);`
	_, err := Conn.Exec(context.Background(), req)
	if err != nil {
		logrus.Error(err)
	}
	return true
}

// SaveDataToDB - saving Metrics to DB
// If metrics already exist on db then update values in DB
func SaveDataToDB(sm *InMemory) error {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return errors.New("error nil Conn")
	}
	sm.Lock()
	defer sm.Unlock()
	// TRUNCATE TABLE COMPANY
	// metrics := make([]types.Metrics, len(sm.MapCounter)+len(sm.MapGauge))
	for k, v := range sm.MapCounter {
		insertCounter := `INSERT INTO praktikum(ID, Type, Delta) values($1, $2, $3)`
		_, err := Conn.Exec(context.Background(), insertCounter, k, "counter", int(v))
		if err != nil {
			updateCounter := `UPDATE praktikum SET Type = $1, Delta = $2 WHERE ID = $3;`
			_, err := Conn.Exec(context.Background(), updateCounter, "counter", int(v), k)
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	for k, v := range sm.MapGauge {
		insertGauge := `INSERT INTO praktikum(ID, Type, Value) values($1, $2, $3)`
		_, err := Conn.Exec(context.Background(), insertGauge, k, "gauge", float64(v))
		if err != nil {
			updateGauge := `UPDATE praktikum SET Type = $1, Value = $2 WHERE ID = $3;`
			_, err := Conn.Exec(context.Background(), updateGauge, "gauge", float64(v), k)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	return nil
}

// RestoreDataFromDB - restoring Metrics from DB
// need call where flag Restore = true
func RestoreDataFromDB(sm *InMemory) error {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return errors.New("error nil Conn")
	}
	sm.Lock()
	defer sm.Unlock()
	ctx := context.Background()
	listCounter := "SELECT ID, Delta FROM praktikum WHERE Type='counter';"
	rowsC, err := Conn.Query(ctx, listCounter)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer rowsC.Close()
	for rowsC.Next() {
		var id string
		var delta int64
		err := rowsC.Scan(&id, &delta)
		if err == nil {
			sm.MapCounter[id] = types.Counter(delta)
		} else {
			logrus.Errorf("error scanning drom db: %v", err)
			return err
		}
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
	defer rowsG.Close()
	for rowsG.Next() {
		var id string
		var value float64
		err := rowsG.Scan(&id, &value)
		if err == nil {
			sm.MapGauge[id] = types.Gauge(value)
		} else {
			logrus.Errorf("error scanning drom db: %v", err)
		}
	}
	err = rowsG.Err()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
