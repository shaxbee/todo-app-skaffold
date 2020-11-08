package routes

import (
	"time"
)

type Opt func(c *config)

func Verbose(verbose bool) Opt {
	return func(c *config) {
		c.Verbose = true
	}
}

func CorsOrigin(origin string) Opt {
	return func(c *config) {
		c.CorsOrigin = origin
	}
}

func CorsRequestHeaders(headers []string) Opt {
	return func(c *config) {
		c.CorsRequestHeaders = headers
	}
}

func CorsMaxAge(maxAge time.Duration) Opt {
	return func(c *config) {
		c.CorsMaxAge = maxAge
	}
}

type config struct {
	Verbose            bool
	CorsOrigin         string
	CorsRequestHeaders []string
	CorsMaxAge         time.Duration
}

var defaultConfig = config{
	CorsOrigin:         "*",
	CorsRequestHeaders: []string{"*"},
}
