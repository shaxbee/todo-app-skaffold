package dbtest

import "os"

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
	migrations string
}

func (c *config) dsn() string {
	return os.Getenv("TEST_DATABASE")
}

var defaultConfig = config{
	image:    "postgres",
	tag:      "13",
	database: "postgres",
	user:     "postgres",
}
