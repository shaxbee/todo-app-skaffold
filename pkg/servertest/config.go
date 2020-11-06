package servertest

import (
	"time"

	"github.com/cenkalti/backoff"
)

type config struct {
	Addr            string
	InitialInterval time.Duration
	MaxElapsedTime  time.Duration
}

func makeConfig(opts ...ConfigOpt) config {
	c := defaultConfig

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

func (c config) ExponentialBackoff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = c.InitialInterval
	bo.MaxElapsedTime = c.MaxElapsedTime

	return bo
}

var defaultConfig = config{
	Addr:            ":http",
	InitialInterval: 10 * time.Millisecond,
	MaxElapsedTime:  5 * time.Second,
}

type ConfigOpt func(*config)

func Addr(addr string) ConfigOpt {
	return func(c *config) {
		c.Addr = addr
	}
}

func InitialInterval(interval time.Duration) ConfigOpt {
	return func(c *config) {
		c.InitialInterval = interval
	}
}

func MaxElapsedTime(elapsedTime time.Duration) ConfigOpt {
	return func(c *config) {
		c.MaxElapsedTime = elapsedTime
	}
}
