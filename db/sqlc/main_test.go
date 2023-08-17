package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {

	testDB, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
