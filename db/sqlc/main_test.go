package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://dev:devpass@localhost:5435/go-simple-bank-v2?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(t *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(t.Run())
}
