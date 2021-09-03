package httprouter

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Debug   string `json:"debug,omitempty"`
}

func DefaultErrorHandler(logger *zap.Logger, verbose bool) ErrorHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		httpErr := AsError(err)
		cause := httpErr.Cause

		if !httpErr.Operational {
			fields := []zap.Field{
				zap.String("path", req.URL.Path),
				zap.Int("status", httpErr.Status),
				zap.String("message", httpErr.Message),
			}

			if cause != nil {
				fields = append(fields, zap.Error(cause))
			}

			logger.Info("http error", fields...)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpErr.Status)

		var debug string
		if verbose && cause != nil {
			debug = cause.Error()
		}

		resp := ErrorResponse{
			Message: httpErr.Message,
			Debug:   debug,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("encode error response", zap.Error(err))
		}
	}
}

func DefaultPanicHandler(logger *zap.Logger, verbose bool) PanicHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, pv interface{}) {
		fields := []zap.Field{
			zap.String("path", req.URL.Path),
		}

		if pv != nil {
			fields = append(fields, zap.Any("panic", pv))
		}

		if verbose {
			fields = append(fields, zap.StackSkip("stack", 1))
		}

		logger.Error("http server panic", fields...)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		resp := ErrorResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("encode error response", zap.Error(err))
		}
	}
}
