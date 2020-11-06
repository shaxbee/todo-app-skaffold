package dbtest

type config struct {
	enabled    bool
	retain     bool
	dsn        string
	image      string
	tag        string
	database   string
	user       string
	migrations string
}

var defaultConfig = config{
	enabled:  true,
	image:    "postgres",
	tag:      "13",
	database: "postgres",
	user:     "postgres",
}

type ConfigOpt func(*config)

func Enabled(enabled bool) ConfigOpt {
	return func(c *config) {
		c.enabled = enabled
	}
}

func Retain(retain bool) ConfigOpt {
	return func(c *config) {
		c.retain = retain
	}
}

func Image(image string) ConfigOpt {
	return func(c *config) {
		c.image = image
	}
}

func DSN(dsn string) ConfigOpt {
	return func(c *config) {
		c.dsn = dsn
	}
}

func Tag(tag string) ConfigOpt {
	return func(c *config) {
		c.tag = tag
	}
}

func Database(database string) ConfigOpt {
	return func(c *config) {
		c.database = database
	}
}

func User(user string) ConfigOpt {
	return func(c *config) {
		c.user = user
	}
}

func Migration(dir string) ConfigOpt {
	return func(c *config) {
		c.migrations = dir
	}
}
