package httprouter

import "net/http"

type ResponseWriter struct {
	delegate   http.ResponseWriter
	statusCode int
}

func NewResponseWriter(delegate http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		delegate: delegate,
	}
}

func (r *ResponseWriter) Header() http.Header {
	return r.delegate.Header()
}

func (r *ResponseWriter) Write(data []byte) (int, error) {
	return r.delegate.Write(data)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.delegate.WriteHeader(statusCode)
}

func (r *ResponseWriter) StatusCode() int {
	return r.statusCode
}
