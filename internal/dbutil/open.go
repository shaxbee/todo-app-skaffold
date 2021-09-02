package dbutil

import (
	"context"
	"database/sql"

	"github.com/cenkalti/backoff/v3"
)

func Open(ctx context.Context, driver, dsn string, opts ...Opt) (*sql.DB, error) {
	c := defaultConfig

	for _, opt := range opts {
		opt(&c)
	}

	var db *sql.DB

	err := backoff.Retry(func() error {
		var err error
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			return err
		}

		return db.Ping()
	}, c.Backoff(ctx))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)

	return db, nil
}
