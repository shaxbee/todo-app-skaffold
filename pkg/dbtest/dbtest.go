package dbtest

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/pressly/goose"
	"github.com/shaxbee/todo-app-skaffold/pkg/dbutil"
)

func Postgres(t *testing.T, opts ...ConfigOpt) *sql.DB {
	t.Helper()

	c := defaultConfig

	for _, opt := range opts {
		opt(&c)
	}

	var dsn = c.dsn

	if c.enabled {
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

		dsn = fmt.Sprintf("port=%s user=%s dbname=%s sslmode=disable", resource.GetPort("5432/tcp"), c.user, c.database)
	}

	db, err := dbutil.Open(
		"postgres",
		dsn,
		dbutil.MaxInterval(100*time.Millisecond),
		dbutil.MaxElapsedTime(10*time.Second),
	)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	t.Logf("connected to database %q", dsn)

	if c.migrations != "" {
		t.Logf("goose: running migrations")

		if err := goose.Up(db, c.migrations); err != nil {
			t.Fatalf("failed to migrate: %v", err)
		}
	}

	return db
}
