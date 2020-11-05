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

type configOpt func(*config)

func Enabled(enabled bool) configOpt {
	return func(c *config) {
		c.enabled = enabled
	}
}

func Retain(retain bool) configOpt {
	return func(c *config) {
		c.retain = retain
	}
}

func Image(image string) configOpt {
	return func(c *config) {
		c.image = image
	}
}

func DSN(dsn string) configOpt {
	return func(c *config) {
		c.dsn = dsn
	}
}

func Tag(tag string) configOpt {
	return func(c *config) {
		c.tag = tag
	}
}

func Database(database string) configOpt {
	return func(c *config) {
		c.database = database
	}
}

func User(user string) configOpt {
	return func(c *config) {
		c.user = user
	}
}

func Migration(dir string) configOpt {
	return func(c *config) {
		c.migrations = dir
	}
}
