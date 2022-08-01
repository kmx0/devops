package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/jackc/pgx/v4"
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

// const (
// 	addr = ":8080" // адрес сервера
// )

func BenchmarkPingDB(b *testing.B) {

	b.Run("PingDB: Before profiling ", func(b *testing.B) {
		ctx := context.Background()
		urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
		for i := 0; i < b.N; i++ {
			PingDBBeforeProfiling(ctx, urlExample)
		}
	})

	b.Run("PingDB: After profiling", func(b *testing.B) {

		ctx := context.Background()
		urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
		for i := 0; i < b.N; i++ {
			PingDBProfiled(ctx, urlExample)
		}
	})

	// go BenchmarkPingDB(b)
	// http.ListenAndServe(addr, nil)
}

func BenchmarkSaveDatatoDB(b *testing.B) {
	ctx := context.Background()
	urlExample := "postgres://postgres:postgres@localhost:5432/metrics"
	PingDBBeforeProfiling(ctx, urlExample)
	b.ResetTimer()
	sm := NewInMemory(config.Config{})
	sm.MapCounter["1"] = types.Counter(1)
	sm.MapCounter["2"] = types.Counter(2)
	sm.MapCounter["3"] = types.Counter(3)
	sm.MapCounter["4"] = types.Counter(4)

	sm.MapGauge["1"] = types.Gauge(1)
	sm.MapGauge["2"] = types.Gauge(2)
	sm.MapGauge["3"] = types.Gauge(3)
	sm.MapGauge["4"] = types.Gauge(4)
	b.Run("SaveDataToDB: Before profiling", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			SaveDataToDBBeforeProfiling(sm)
		}
	})

	b.Run("SaveDataToDB: After profiling", func(b *testing.B) {

		sm := NewInMemory(config.Config{})

		for i := 0; i < b.N; i++ {
			SaveDataToDBProfiled(sm)
		}
	})

	// go BenchmarkPingDB(b)
	// http.ListenAndServe(addr, nil)
}

func PingDBBeforeProfiling(ctx context.Context, urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"

	var err error

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	Conn, err = pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
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

func PingDBProfiled(ctx context.Context, urlExample string) bool {
	var err error

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	if Conn != nil {
		return true
	}
	Conn, err = pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
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

func SaveDataToDBBeforeProfiling(sm *InMemory) {
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

func SaveDataToDBProfiled(sm *InMemory) {
	if Conn == nil {
		logrus.Error("Error nil Conn")
		return
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

}

func TestPingDB(t *testing.T) {

	type wantStruct struct {
		res bool
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Fail",
			want: wantStruct{
				res: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := PingDB(context.Background(), "")
			// PingDB(ctx context.Context, urlExample string) bool

			assert.Equal(t, tt.want.res, res)

		})
	}
}

func TestCheckTableExist(t *testing.T) {

	type wantStruct struct {
		res bool
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Fail",
			want: wantStruct{
				res: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CheckTableExist()
			// PingDB(ctx context.Context, urlExample string) bool

			assert.Equal(t, tt.want.res, res)

		})
	}
}
func TestAddTabletoDB(t *testing.T) {

	type wantStruct struct {
		res bool
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Fail",
			want: wantStruct{
				res: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := AddTabletoDB()
			// PingDB(ctx context.Context, urlExample string) bool

			assert.Equal(t, tt.want.res, res)

		})
	}
}

func TestSaveDataToDB(t *testing.T) {

	cfg := config.Config{}
	sm := NewInMemory(cfg)
	type wantStruct struct {
		err error
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Fail",
			want: wantStruct{
				err: errors.New("error nil Conn"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveDataToDB(sm)
			assert.Equal(t, tt.want.err.Error(), err.Error())
		})
	}
}

func TestRestoreDataFromDB(t *testing.T) {
	cfg := config.Config{}
	sm := NewInMemory(cfg)
	type wantStruct struct {
		err error
	}

	tests := []struct {
		name string
		want wantStruct
	}{
		{
			name: "Fail",
			want: wantStruct{
				err: errors.New("error nil Conn"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RestoreDataFromDB(sm)
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}
		})
	}
}
