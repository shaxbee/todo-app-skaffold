package cors

import (
	"net/http"

	"github.com/shaxbee/todo-app-skaffold/internal/httprouter"
)

func RouterOpts(opts ...Opt) []httprouter.Opt {
	return []httprouter.Opt{
		httprouter.WithGlobalOptions(GlobalOptions(opts...)),
		httprouter.WithMiddleware(Middleware(opts...)),
	}
}

func GlobalOptions(opts ...Opt) httprouter.HandlerFunc {
	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}

	return func(w http.ResponseWriter, req *http.Request) error {
		if req.Header.Get("Access-Control-Request-Method") == "" {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}

		header := w.Header()

		origin := req.Header.Get("Origin")
		requestHeaders := req.Header.Get("Access-Control-Request-Headers")

		switch {
		case c.OriginWildcard() && origin != "":
			header.Set("Access-Control-Allow-Origin", origin)
			header.Set("Vary", origin)
		case c.OriginWildcard():
			header.Set("Access-Control-Allow-Origin", c.origin)
		default:
			return nil
		}

		switch {
		case c.RequestHeadersWildcard() && requestHeaders != "":
			header.Set("Access-Control-Allow-Headers", requestHeaders)
		case requestHeaders != "":
			header.Set("Access-Control-Allow-Headers", c.requestHeaders)
		}

		header.Set("Access-Control-Allow-Methods", header.Get("Allow"))

		if c.allowCredentials {
			header.Set("Access-Control-Allow-Credentials", "true")
		}

		if c.maxAge != "" {
			header.Set("Access-Control-Max-Age", c.maxAge)
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}

func Middleware(opts ...Opt) httprouter.Middleware {
	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}

	return func(next httprouter.HandlerFunc) httprouter.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) error {
			header := w.Header()

			header.Set("Access-Control-Allow-Origin", c.origin)

			switch {
			case c.RequestHeadersWildcard() && c.allowCredentials:
				header.Set("Access-Control-Allow-Headers", "*, Authorization")
			case c.RequestHeadersWildcard():
				header.Set("Access-Control-Allow-Headers", "*")
			default:
				header.Set("Access-Control-Allow-Headers", c.requestHeaders)
			}

			return next(w, req)
		}
	}
}
