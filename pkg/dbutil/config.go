package dbutil

import (
	"time"

	"github.com/cenkalti/backoff"
)

type config struct {
	MaxInterval    time.Duration
	MaxElapsedTime time.Duration
	MaxIdleConns   int
	MaxOpenConns   int
}

func (c config) ExponentialBackOff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = c.MaxInterval
	bo.MaxElapsedTime = c.MaxElapsedTime

	return bo
}

var defaultConfig = config{
	MaxInterval:    5 * time.Second,
	MaxElapsedTime: 1 * time.Minute,
	MaxIdleConns:   2,
	MaxOpenConns:   0,
}

type configOpt func(*config)

func MaxInterval(interval time.Duration) configOpt {
	return func(c *config) {
		c.MaxInterval = interval
	}
}

func MaxElapsedTime(elapsedTime time.Duration) configOpt {
	return func(c *config) {
		c.MaxElapsedTime = elapsedTime
	}
}

func MaxIdleConns(idleConns int) configOpt {
	return func(c *config) {
		c.MaxIdleConns = idleConns
	}
}

func MaxOpenConns(openConns int) configOpt {
	return func(c *config) {
		c.MaxOpenConns = openConns
	}
}
