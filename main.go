package main

import (
	"database/sql"
	"log"

	"github.com/freer4an/simple-bank/api"
	db "github.com/freer4an/simple-bank/db/sqlc"
	"github.com/freer4an/simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.InitConfig(".")
	if err != nil {
		log.Fatal("error reading config", err)
	}

	conn, err := sql.Open(config.DB_driver, config.DB_source)
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(config.SB_ADDR); err != nil {
		log.Fatal(err)
	}
}
