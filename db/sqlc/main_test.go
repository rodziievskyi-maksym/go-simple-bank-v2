package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rodziievskyi-maksym/go-simple-bank-v2/config"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(t *testing.M) {
	cfg, err := config.Load("../../")
	if err != nil {
		log.Fatal("cannot load config file:", err)
	}

	testDB, err = sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(t.Run())
}
