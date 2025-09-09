package main

import (
	"database/sql"
	"log"

	"github.com/rodziievskyi-maksym/go-simple-bank-v2/api"
	"github.com/rodziievskyi-maksym/go-simple-bank-v2/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://dev:devpass@localhost:5435/go-simple-bank-v2?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.ServeHTTP(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
