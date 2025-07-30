package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/0xOnah/bank/internal/config"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	cfg, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(dbDriver, cfg.DSN)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
