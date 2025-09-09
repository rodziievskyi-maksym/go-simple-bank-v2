package main

import (
	"database/sql"
	"log"

	"github.com/rodziievskyi-maksym/go-simple-bank-v2/api"
	"github.com/rodziievskyi-maksym/go-simple-bank-v2/config"
	"github.com/rodziievskyi-maksym/go-simple-bank-v2/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal("cannot load cfg file: ", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.ServeHTTP(cfg.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
