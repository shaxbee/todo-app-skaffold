package routes

import (
	"strings"
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
		c.CorsRequestHeaders = strings.Join(headers, ",")
	}
}

func CorsAllowCredentials(allowCredentials bool) Opt {
	return func(c *config) {
		c.CorsAllowCredentials = allowCredentials
	}
}

func CorsMaxAge(maxAge time.Duration) Opt {
	return func(c *config) {
		c.CorsMaxAge = maxAge
	}
}

type config struct {
	Verbose              bool
	CorsOrigin           string
	CorsRequestHeaders   string
	CorsAllowCredentials bool
	CorsMaxAge           time.Duration
}

var defaultConfig = config{
	CorsOrigin:           "*",
	CorsRequestHeaders:   "*",
	CorsAllowCredentials: true,
}
