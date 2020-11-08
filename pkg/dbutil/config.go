package dbutil

import (
	"time"

	"github.com/cenkalti/backoff/v3"
)

type config struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func (c config) ExponentialBackOff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = c.InitialInterval
	bo.MaxInterval = c.MaxInterval
	bo.MaxElapsedTime = c.MaxElapsedTime

	return bo
}

var defaultConfig = config{
	InitialInterval: 100 * time.Millisecond,
	MaxInterval:     5 * time.Second,
	MaxElapsedTime:  1 * time.Minute,
	MaxIdleConns:    2,
	MaxOpenConns:    0,
}

type ConfigOpt func(*config)

func InitialInterval(interval time.Duration) ConfigOpt {
	return func(c *config) {
		c.InitialInterval = interval
	}
}

func MaxInterval(interval time.Duration) ConfigOpt {
	return func(c *config) {
		c.MaxInterval = interval
	}
}

func MaxElapsedTime(elapsedTime time.Duration) ConfigOpt {
	return func(c *config) {
		c.MaxElapsedTime = elapsedTime
	}
}

func MaxIdleConns(idleConns int) ConfigOpt {
	return func(c *config) {
		c.MaxIdleConns = idleConns
	}
}

func MaxOpenConns(openConns int) ConfigOpt {
	return func(c *config) {
		c.MaxOpenConns = openConns
	}
}
