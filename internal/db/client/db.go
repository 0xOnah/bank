package client

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type DBClient struct {
	Client *sql.DB
}

func NewDBClient(DSN string) (*DBClient, error) {
	conn, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	return &DBClient{Client: conn}, nil
}

func (dc *DBClient) Ping(ctx context.Context) error {
	return dc.Client.PingContext(ctx)
}

func (dc *DBClient) Migrate() {

}
