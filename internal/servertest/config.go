package servertest

import (
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v3"
)

type Opt func(*config)

type MakeHandlerFunc func() http.Handler

func Backoff(bo backoff.BackOff) Opt {
	return func(c *config) {
		c.backoff = bo
	}
}

func MakeHandler(makeHandler MakeHandlerFunc) Opt {
	return func(c *config) {
		c.makeHandler = makeHandler
	}
}

type config struct {
	backoff     backoff.BackOff
	makeHandler MakeHandlerFunc
}

func (c config) Endpoint() string {
	endpoint := os.Getenv("TEST_ENDPOINT")
	switch {
	case endpoint != "":
		return endpoint
	case endpoint == "" && c.makeHandler == nil:
		return "http://:80"
	default:
		return ""
	}
}

func (c config) Backoff() backoff.BackOff {
	if c.backoff != nil {
		return c.backoff
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 10 * time.Millisecond
	bo.MaxElapsedTime = 5 * time.Second

	return bo
}
