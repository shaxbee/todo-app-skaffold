package dbutil

import (
	"database/sql"

	"github.com/cenkalti/backoff"
)

func Open(driver string, dsn string, opts ...configOpt) (*sql.DB, error) {
	c := defaultConfig

	for _, opt := range opts {
		opt(&c)
	}

	var db *sql.DB

	err := backoff.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return err
		}

		return db.Ping()
	}, c.ExponentialBackOff())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)

	return db, nil
}
