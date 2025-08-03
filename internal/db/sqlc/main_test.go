package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, os.Getenv("DSN"))
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	driver, err := postgres.WithInstance(testDB, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to create driver db", err)
	}
	ma, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err != nil {
		log.Fatal("failed to create migrate instance", err)
	}
	ma.Up()
	testQueries = New(testDB)

	os.Exit(m.Run())
}
