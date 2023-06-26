package cockroachdb

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
)

type CockroachDB struct {
	db  *sql.DB
	ctx context.Context
	log log.Logger
}

func Connect() *CockroachDB {
	dsn := "postgresql://shutt:" + os.Getenv("COCKROACH_DB_PASS") + "@tft-meta-8562.8nj.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	var now time.Time
	err = db.QueryRow("SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	return &CockroachDB{db: db, ctx: context.Background(), log: *log.New(log.Writer(), "db log", 755)}
}
