package dbtest

import (
	"os"
	"time"

	"github.com/cenkalti/backoff/v3"
)

type Opt func(*config)

func Retain(retain bool) Opt {
	return func(c *config) {
		c.retain = retain
	}
}

func Image(image string) Opt {
	return func(c *config) {
		c.image = image
	}
}

func Tag(tag string) Opt {
	return func(c *config) {
		c.tag = tag
	}
}

func Database(database string) Opt {
	return func(c *config) {
		c.database = database
	}
}

func User(user string) Opt {
	return func(c *config) {
		c.user = user
	}
}

func Backoff(bo backoff.BackOff) Opt {
	return func(c *config) {
		c.backoff = bo
	}
}

func Migration(dir string) Opt {
	return func(c *config) {
		c.migrations = dir
	}
}

type config struct {
	retain     bool
	image      string
	tag        string
	database   string
	user       string
	backoff    backoff.BackOff
	migrations string
}

func (c *config) DSN() string {
	return os.Getenv("TEST_DATABASE")
}

func (c *config) Backoff() backoff.BackOff {
	if c.backoff != nil {
		return c.backoff
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 10 * time.Millisecond
	bo.MaxInterval = 1 * time.Second
	bo.MaxElapsedTime = 10 * time.Second

	return bo
}

var defaultConfig = config{
	image:    "postgres",
	tag:      "13",
	database: "postgres",
	user:     "postgres",
}
