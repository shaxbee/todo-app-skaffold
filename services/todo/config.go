package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Dev    bool `json:"dev" envconfig:"DEV" desc:"Development mode"`
	Server struct {
		Addr            string        `json:"public_addr" envconfig:"Addr" default:":80" desc:"Server listen address"`
		Timeout         time.Duration `json:"timeout" envconfig:"TIMEOUT" default:"5s" desc:"Operation timeout"`
		ShutdownTimeout time.Duration `json:"shutdown_timeout" envconfig:"SHUTDOWN_TIMEOUT" default:"10s" desc:"Shutdown timeout"`
	} `json:"server" envconfig:"SERVER"`
	DB struct {
		DSN          string `json:"dsn" envconfig:"DSN" default:"postgres://" desc:"Database data source name"`
		MaxIdleConns int    `json:"max_idle_conns" envconfig:"MAX_IDLE_CONNS" default:"5" desc:"Database max idle connections"`
		MaxOpenConns int    `json:"max_open_conns" envconfig:"MAX_OPEN_CONNS" default:"20" desc:"Database max open connections"`
	} `json:"db" envconfig:"DB"`
}

func parseConfig() (*config, error) {
	var c config
	if err := envconfig.Process("TODO", &c); err != nil {
		return nil, err
	}

	return &c, nil
}
