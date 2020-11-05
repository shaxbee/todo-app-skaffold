package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/errgroup"

	"github.com/shaxbee/todo-app-skaffold/pkg/dbutil"
	"github.com/shaxbee/todo-app-skaffold/pkg/httperror"
)

type Container struct {
	config     *Config
	db         *sql.DB
	todoServer *TodoServer
	httpServer *http.Server
	listener   net.Listener

	once struct {
		db, todoServer, listener, httpServer sync.Once
	}
}

func NewContainer(config *Config, opts ...containerOpt) *Container {
	c := Container{
		config: config,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *Container) DB() *sql.DB {
	c.once.db.Do(func() {
		if c.db != nil {
			return
		}

		var err error
		c.db, err = dbutil.Open(
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

func (c *Container) TodoServer() *TodoServer {
	c.once.todoServer.Do(func() {
		c.todoServer = New(c.DB())
	})

	return c.todoServer
}

func (c *Container) HTTPServer() *http.Server {
	c.once.httpServer.Do(func() {
		todoServer := c.TodoServer()

		router := httprouter.New()
		errorMiddleware := httperror.NewMiddleware(httperror.Verbose(c.config.Dev))
		todoServer.RegisterRoutes(router, errorMiddleware)

		c.httpServer = &http.Server{
			Addr:         c.config.Server.Addr,
			ReadTimeout:  c.config.Server.Timeout,
			WriteTimeout: c.config.Server.Timeout,
			Handler:      router,
		}
	})

	return c.httpServer
}

func (c *Container) Listener() net.Listener {
	c.once.listener.Do(func() {
		var err error
		c.listener, err = net.Listen("tcp", c.config.Server.Addr)
		if err != nil {
			log.Fatalf("failed to listen at %q: %v", c.config.Server.Addr, err)
		}
	})

	return c.listener
}

func (c *Container) Run(ctx context.Context) *errgroup.Group {
	server := c.HTTPServer()

	log.Println("server created")

	listener := c.Listener()
	addr := listener.Addr()
	log.Printf("listening at %q", addr)

	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), c.config.Server.ShutdownTimeout)
		defer cancel()

		log.Printf("shutdown initiated")

		if err := server.Shutdown(sctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		log.Printf("server listening at %q has shutdown", addr)
		return nil
	})

	errg.Go(func() error {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	return errg
}

type containerOpt func(c *Container)

func ContainerDB(db *sql.DB) containerOpt {
	return func(c *Container) {
		c.db = db
	}
}
