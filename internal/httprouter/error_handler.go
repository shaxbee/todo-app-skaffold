package httprouter

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Debug   string `json:"debug,omitempty"`
}

func DefaultErrorHandler(w http.ResponseWriter, req *http.Request, verbose bool, err error) {
	httpErr := AsError(err)

	var debug string
	switch {
	case verbose && httpErr.Cause != nil:
		debug = httpErr.Cause.Error()
		log.Printf("http error %d %s: %+v", httpErr.Status, httpErr.Message, httpErr.Cause)
	case verbose:
		log.Printf("http error %d %s", httpErr.Status, httpErr.Message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)

	data, err := json.Marshal(ErrorResponse{
		Message: httpErr.Message,
		Debug:   debug,
	})
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	_, _ = w.Write(data)
}
