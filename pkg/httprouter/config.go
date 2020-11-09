package httprouter

import (
	"strconv"
	"strings"
	"time"
)

type Opt func(c *config)

func Verbose(verbose bool) Opt {
	return func(c *config) {
		c.verbose = true
	}
}

func CorsEnabled(enabled bool) Opt {
	return func(c *config) {
		c.corsEnabled = enabled
	}
}

func CorsOrigin(origin string) Opt {
	return func(c *config) {
		c.corsOrigin = origin
	}
}

func CorsRequestHeaders(headers []string) Opt {
	return func(c *config) {
		c.corsRequestHeaders = strings.Join(headers, ", ")
	}
}

func CorsAllowCredentials(allowCredentials bool) Opt {
	return func(c *config) {
		c.corsAllowCredentials = allowCredentials
	}
}

func CorsMaxAge(maxAge time.Duration) Opt {
	return func(c *config) {
		c.corsMaxAge = strconv.FormatInt(int64(maxAge), 10)
	}
}

type config struct {
	verbose              bool
	corsEnabled          bool
	corsOrigin           string
	corsRequestHeaders   string
	corsAllowCredentials bool
	corsMaxAge           string
}

func (c config) CorsOriginWildcard() bool {
	return c.corsOrigin == "*"
}

func (c config) CorsRequestHeadersWildcard() bool {
	return c.corsRequestHeaders == "*"
}

var defaultConfig = config{
	corsOrigin:           "*",
	corsRequestHeaders:   "*",
	corsAllowCredentials: true,
}
