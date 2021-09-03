package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/shaxbee/todo-app-skaffold/internal/httprouter"
	"github.com/shaxbee/todo-app-skaffold/internal/middleware/cors"
	"github.com/shaxbee/todo-app-skaffold/services/todo"
)

type container struct {
	config *Config

	state struct {
		logger     *zap.Logger
		db         *sql.DB
		todoServer *todo.Server
		httpRouter *httprouter.Router
		httpServer *http.Server
		listener   net.Listener
	}

	once struct {
		logger, db, todoServer, httpRouter, httpServer, listener sync.Once
	}
}

func newContainer(config *Config) *container {
	return &container{
		config: config,
	}
}

func (c *container) logger() *zap.Logger {
	c.once.logger.Do(func() {
		var lc zap.Config
		switch {
		case c.config.Dev:
			lc = zap.NewDevelopmentConfig()
		default:
			lc = zap.NewProductionConfig()
		}

		logger, err := lc.Build()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		c.state.logger = logger
	})

	return c.state.logger
}

func (c *container) db() *sql.DB {
	c.once.db.Do(func() {
		if c.state.db != nil {
			return
		}

		db, err := sql.Open("pgx", c.config.DB.DSN)
		if err != nil {
			c.logger().Fatal("db", zap.String("dsn", c.config.DB.DSN), zap.Error(err))
		}

		db.SetMaxIdleConns(c.config.DB.MaxIdleConns)
		db.SetMaxOpenConns(c.config.DB.MaxOpenConns)

		c.state.db = db
	})

	return c.state.db
}

func (c *container) todoServer() *todo.Server {
	c.once.todoServer.Do(func() {
		c.state.todoServer = todo.NewServer(c.db())
	})

	return c.state.todoServer
}

func (c *container) httpRouter() *httprouter.Router {
	c.once.httpRouter.Do(func() {
		todoServer := c.todoServer()

		opts := []httprouter.Opt{
			httprouter.WithVerbose(c.config.Dev),
		}

		if c.config.Server.CorsEnabled {
			opts = append(opts, cors.RouterOpts()...)
		}

		router := httprouter.New(c.logger(), opts...)
		todoServer.RegisterRoutes(router)

		c.state.httpRouter = router
	})

	return c.state.httpRouter
}

func (c *container) httpServer() *http.Server {
	c.once.httpServer.Do(func() {
		c.state.httpServer = &http.Server{
			Addr:         c.config.Server.Addr,
			ReadTimeout:  c.config.Server.Timeout,
			WriteTimeout: c.config.Server.Timeout,
			Handler:      c.httpRouter(),
		}
	})

	return c.state.httpServer
}

func (c *container) listener() net.Listener {
	c.once.listener.Do(func() {
		listener, err := net.Listen("tcp", c.config.Server.Addr)
		if err != nil {
			c.logger().Fatal("listener", zap.String("addr", c.config.Server.Addr), zap.Error(err))
		}

		c.state.listener = listener
	})

	return c.state.listener
}
