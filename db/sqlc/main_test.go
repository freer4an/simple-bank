package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/freer4an/simple-bank/util"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {

	config, err := util.InitConfig("../..")
	if err != nil {
		log.Fatal("error reading config", err)
	}

	testDB, err := sql.Open(config.DB_driver, config.DB_source)
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
