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

var testStore *Store

func TestMain(m *testing.M) {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	testDB, err := sql.Open(driverName, os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
