// +build integration

package dbtest

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"testing"

	"github.com/cenkalti/backoff/v3"
	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose"
	"github.com/shaxbee/todo-app-skaffold/pkg/dbutil"
)

func SetupPostgres(t testing.TB, opts ...Opt) *sql.DB {
	t.Helper()

	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}

	if dsn := c.DSN(); dsn != "" {
		return openDB(t, dsn, c.migrations, c.Backoff())
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	name := containerName("postgres")

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       name,
		Repository: c.database,
		Tag:        c.tag,
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", c.user),
			fmt.Sprintf("POSTGRES_DB=%s", c.database),
			"POSTGRES_HOST_AUTH_METHOD=trust",
		},
	})
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}

	t.Cleanup(func() {
		if t.Failed() && c.retain {
			return
		}

		if err := pool.Purge(resource); err != nil {
			t.Errorf("failed to purge pool: %v", err)
		}
	})

	t.Logf("started container %q", name)

	dsn := fmt.Sprintf("port=%s user=%s dbname=%s sslmode=disable", resource.GetPort("5432/tcp"), c.user, c.database)
	return openDB(t, dsn, c.migrations, c.Backoff())
}

func openDB(t testing.TB, dsn, migrations string, bo backoff.BackOff) *sql.DB {
	db, err := dbutil.Open(context.Background(), "postgres", dsn, dbutil.Backoff(bo))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Logf("connected to database %q", dsn)

	if migrations != "" {
		t.Logf("goose: running migrations")

		if err := goose.Up(db, migrations); err != nil {
			t.Fatalf("failed to migrate: %v", err)
		}
	}

	return db
}

func containerName(driver string) string {
	suffix, _ := rand.Int(rand.Reader, big.NewInt(100000))
	return fmt.Sprintf("%s-%s", driver, suffix)
}
