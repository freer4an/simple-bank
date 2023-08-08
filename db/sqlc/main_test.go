package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const driverName = "postgres"
const db_source = "postgresql://ady:password@localhost:5432/simple_bank?sslmode=disable"

var testStore *Store

func TestMain(m *testing.M) {
	var err error

	// if err = godotenv.Load(); err != nil {
	// 	log.Fatal(err)
	// }

	testDB, err := sql.Open(driverName, db_source)
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
