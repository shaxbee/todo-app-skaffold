package servertest

import (
	"net"
	"testing"

	"github.com/cenkalti/backoff/v3"
)

func WaitForServer(t *testing.T, opts ...ConfigOpt) {
	t.Helper()

	c := makeConfig(opts...)

	err := backoff.Retry(func() error {
		conn, err := net.Dial("tcp", c.Addr)
		if err != nil {
			return err
		}
		conn.Close()

		return nil
	}, c.ExponentialBackoff())
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
}
