package httprouter

import "net/http"

type Opt func(c *config)

func Verbose(verbose bool) Opt {
	return func(c *config) {
		c.verbose = true
	}
}

func ErrorHandler(handler ErrorHandlerFunc) Opt {
	return func(c *config) {
		c.errorHandler = handler
	}
}

func Middleware(middleware ...MiddlewareFunc) Opt {
	return func(c *config) {
		c.middleware = append(c.middleware, middleware...)
	}
}

func GlobalOptions(handler http.HandlerFunc) Opt {
	return func(c *config) {
		c.globalOptions = handler
	}
}

type config struct {
	verbose       bool
	errorHandler  ErrorHandlerFunc
	middleware    []MiddlewareFunc
	globalOptions http.HandlerFunc
}

var defaultConfig = config{
	errorHandler: DefaultErrorHandler,
}
