package db

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations
var migrationFs embed.FS

func virtualFS() (fs.FS, error) {
	vs, err := fs.Sub(migrationFs, "migrations")
	if err != nil {
		return nil, err
	}
	return vs, nil
}

type DBClient struct {
	Client *sql.DB
}

func NewDBClient(DSN string) (*DBClient, error) {
	conn, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	//TODO: currently hardcoded handle that later to be passed at runtime
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxIdleTime(time.Second * 15)
	return &DBClient{Client: conn}, nil
}

func (dc *DBClient) Ping(ctx context.Context) error {
	return dc.Client.PingContext(ctx)
}

func (dc *DBClient) MigrateUP() error {
	migrations, err := virtualFS()
	if err != nil {
		return err
	}
	sourceDriver, err := iofs.New(migrations, ".")
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(dc.Client, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
