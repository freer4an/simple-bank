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
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	testDB, err = sql.Open(driverName, os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
