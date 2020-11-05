package dbtest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/pressly/goose"
	"github.com/shaxbee/todo-app-skaffold/pkg/dbutil"
)

var db *sql.DB

func TestMain(m *testing.M, opts ...configOpt) {
	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}

	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
		dsn      = c.dsn
		err      error
	)
	if c.enabled {
		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("failed to create pool: %v", err)
		}

		name := containerName("postgres")
		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
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
			log.Fatalf("failed to start postgres: %v", err)
		}

		log.Printf("started container %q", name)

		dsn = fmt.Sprintf("port=%s user=%s dbname=%s sslmode=disable", resource.GetPort("5432/tcp"), c.user, c.database)
	}

	db, err = dbutil.Open(
		"postgres",
		dsn,
		dbutil.MaxInterval(100*time.Millisecond),
		dbutil.MaxElapsedTime(10*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	log.Printf("connected to database %q", dsn)

	if c.migrations != "" {
		log.Println("goose: running migrations")

		if err := goose.Up(db, c.migrations); err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}

	code := m.Run()

	if pool != nil && code != 0 && !c.retain {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("failed to purge pool: %v", err)
		}
	}

	os.Exit(code)
}

func DB() *sql.DB {
	return db
}
