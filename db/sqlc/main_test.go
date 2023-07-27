package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const driverName = "postgres"

var testQueries *Queries

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(driverName, os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
