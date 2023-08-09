package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/freer4an/simple-bank/api"
	db "github.com/freer4an/simple-bank/db/sqlc"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const driverName = "postgres"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(driverName, os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(os.Getenv("ADDR")); err != nil {
		log.Fatal(err)
	}
}
