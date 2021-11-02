package httprouter

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(delegate http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: delegate,
	}
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseWriter) StatusCode() int {
	return r.statusCode
}
