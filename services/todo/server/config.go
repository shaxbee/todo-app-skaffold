package server

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev    bool `json:"dev" envconfig:"DEV" desc:"Development mode"`
	Server struct {
		Addr            string        `json:"public_addr" envconfig:"ADDR" default:":http" desc:"Server listen address"`
		Timeout         time.Duration `json:"timeout" envconfig:"TIMEOUT" default:"5s" desc:"Operation timeout"`
		ShutdownTimeout time.Duration `json:"shutdown_timeout" envconfig:"SHUTDOWN_TIMEOUT" default:"10s" desc:"Shutdown timeout"`
		CorsEnabled     bool          `json:"cors_enabled" envconfig:"CORS_ENABLED" default:"false" desc:"Enable cors"`
	} `json:"server" envconfig:"SERVER"`
	DB struct {
		DSN          string `json:"dsn" envconfig:"DSN" default:"" desc:"Database data source name"`
		MaxIdleConns int    `json:"max_idle_conns" envconfig:"MAX_IDLE_CONNS" default:"5" desc:"Database max idle connections"`
		MaxOpenConns int    `json:"max_open_conns" envconfig:"MAX_OPEN_CONNS" default:"20" desc:"Database max open connections"`
	} `json:"db" envconfig:"DB"`
}

func ParseConfig(opts ...ConfigOpt) (*Config, error) {
	var c Config

	if err := envconfig.Process("TODO", &c); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c, nil
}

type ConfigOpt func(*Config)

func ConfigAddr(addr string) ConfigOpt {
	return func(c *Config) {
		c.Server.Addr = addr
	}
}
