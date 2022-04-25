package storage

import (
	"context"
	_ "database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func PingDB(urlExample string) bool {
	// urlExample := "postgres://postgres:postgres@localhost:5432/metrics"

	// host := "127.0.0.1"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, urlExample)
	// fmt.Sprintf("host=%s port='5432' dbname='metrics' user='postgres' password='postgres' sslmode=disable", host))
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// os.Exit(1)
		return false
	}
	defer conn.Close(context.Background())
	err = conn.Ping(ctx)
	return err == nil
}
