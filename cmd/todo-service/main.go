package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	config, err := parseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	container := newContainer(config)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, container); err != nil {
		container.logger().Error("run", zap.Error(err))
	}
}

func run(ctx context.Context, c *container) error {
	server := c.httpServer()

	listener := c.listener()
	addr := listener.Addr().String()
	c.logger().Info("server started", zap.String("addr", addr))

	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), c.config.Server.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(sctx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}

		c.logger().Info("server shutdown", zap.String("addr", addr))
		return nil
	})

	errg.Go(func() error {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	return errg.Wait()
}
