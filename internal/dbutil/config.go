package dbutil

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"
)

var defaultConfig = config{
	MaxIdleConns: 2,
	MaxOpenConns: 0,
}

type Opt func(*config)

func Backoff(bo backoff.BackOff) Opt {
	return func(c *config) {
		c.backoff = bo
	}
}

func MaxIdleConns(idleConns int) Opt {
	return func(c *config) {
		c.MaxIdleConns = idleConns
	}
}

func MaxOpenConns(openConns int) Opt {
	return func(c *config) {
		c.MaxOpenConns = openConns
	}
}

type config struct {
	backoff      backoff.BackOff
	MaxIdleConns int
	MaxOpenConns int
}

func (c config) Backoff(ctx context.Context) backoff.BackOff {
	if c.backoff != nil {
		return backoff.WithContext(c.backoff, ctx)
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 100 * time.Millisecond
	bo.MaxInterval = 5 * time.Second
	bo.MaxElapsedTime = 1 * time.Minute

	return backoff.WithContext(bo, ctx)
}
