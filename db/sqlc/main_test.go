package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rafaelvitoadrian/simplebank2/utils"
)

var TestQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../../.")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}

	TestQueries = New(testDB)
	os.Exit(m.Run())
}
