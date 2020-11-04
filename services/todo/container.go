package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/julienschmidt/httprouter"

	"github.com/shaxbee/todo-app-skaffold/pkg/httperror"
	"github.com/shaxbee/todo-app-skaffold/services/todo/server"
)

type container struct {
	config *config
	db     *sql.DB

	once struct {
		db sync.Once
	}
}

func (c *container) DB() *sql.DB {
	c.once.db.Do(func() {
		bo := backoff.NewExponentialBackOff()
		bo.MaxInterval = 5 * time.Second
		bo.MaxElapsedTime = 1 * time.Minute

		err := backoff.Retry(func() error {
			db, err := sql.Open("postgres", c.config.DB.DSN)
			if err != nil {
				return err
			}

			c.db = db
			return nil
		}, bo)
		if err != nil {
			log.Fatal("failed to connect to database: %w", err)
		}

		c.db.SetMaxIdleConns(c.config.DB.MaxIdleConns)
		c.db.SetMaxOpenConns(c.config.DB.MaxOpenConns)
	})

	return c.db
}

func (c *container) TodoServer() *server.TodoServer {
	return server.New(c.DB())
}

func (c *container) HTTPServer() *http.Server {
	todoServer := c.TodoServer()

	router := httprouter.New()
	errorMiddleware := httperror.NewMiddleware(httperror.Verbose(c.config.Dev))
	todoServer.RegisterRoutes(router, errorMiddleware)

	return &http.Server{
		Addr:         c.config.Server.Addr,
		ReadTimeout:  c.config.Server.Timeout,
		WriteTimeout: c.config.Server.Timeout,
		Handler:      router,
	}
}

func (c *container) Run(ctx context.Context) error {
	server := c.HTTPServer()

	go func() {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(ctx, c.config.Server.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(sctx); err != nil {
			log.Printf("failed to shutdown server: %v", err)
		}
	}()

	log.Printf("todo listening at %s", c.config.Server.Addr)
	return server.ListenAndServe()
}
