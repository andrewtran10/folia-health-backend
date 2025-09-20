package dal

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite" // SQLite driver
)

func NewDatabaseConn(ctx context.Context, dbLocation string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbLocation)
	if err != nil {
		return nil, err
	}

	return db, nil
}
