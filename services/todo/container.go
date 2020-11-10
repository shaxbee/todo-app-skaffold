package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/shaxbee/todo-app-skaffold/pkg/dbutil"
	"github.com/shaxbee/todo-app-skaffold/pkg/httprouter"
	"github.com/shaxbee/todo-app-skaffold/services/todo/server"
)

type container struct {
	config     *Config
	db         *sql.DB
	todoServer *server.TodoServer
	httpServer *http.Server
	listener   net.Listener

	once struct {
		db, todoServer, listener, httpServer sync.Once
	}
}

func newContainer(config *Config) *container {
	return &container{
		config: config,
	}
}

func (c *container) DB(ctx context.Context) *sql.DB {
	c.once.db.Do(func() {
		if c.db != nil {
			return
		}

		var err error
		c.db, err = dbutil.Open(
			ctx,
			"postgres",
			c.config.DB.DSN,
			dbutil.MaxIdleConns(c.config.DB.MaxIdleConns),
			dbutil.MaxOpenConns(c.config.DB.MaxOpenConns),
		)
		if err != nil {
			log.Fatal("failed to connect to database: %w", err)
		}
	})

	return c.db
}

func (c *container) TodoServer(ctx context.Context) *server.TodoServer {
	c.once.todoServer.Do(func() {
		c.todoServer = server.New(c.DB(ctx))
	})

	return c.todoServer
}

func (c *container) HTTPServer(ctx context.Context) *http.Server {
	c.once.httpServer.Do(func() {
		todoServer := c.TodoServer(ctx)

		router := httprouter.New(
			httprouter.Verbose(c.config.Dev),
			httprouter.CorsEnabled(c.config.Server.CorsEnabled),
		)

		todoServer.RegisterRoutes(router)

		c.httpServer = &http.Server{
			Addr:         c.config.Server.Addr,
			ReadTimeout:  c.config.Server.Timeout,
			WriteTimeout: c.config.Server.Timeout,
			Handler:      router,
		}
	})

	return c.httpServer
}

func (c *container) Listener() net.Listener {
	c.once.listener.Do(func() {
		var err error
		c.listener, err = net.Listen("tcp", c.config.Server.Addr)
		if err != nil {
			log.Fatalf("failed to listen at %q: %v", c.config.Server.Addr, err)
		}
	})

	return c.listener
}

func (c *container) Addr() string {
	return c.listener.Addr().String()
}

func (c *container) Run(ctx context.Context) *errgroup.Group {
	httpServer := c.HTTPServer(ctx)

	listener := c.Listener()
	log.Printf("listening at %q", c.Addr())

	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), c.config.Server.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(sctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		log.Println("server has shutdown", c.Addr())
		return nil
	})

	errg.Go(func() error {
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	return errg
}
