package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/shaxbee/todo-app-skaffold/pkg/httperror"
)

// JSONRequestBody expects application/json content-type and attempts
// to unmarshal request body into dst.
// 415 Unsupported Media Type is returned if invalid content-type was provided
// 400 Bad Request is returned if request body failed to unmarshal
func JSONRequestBody(req *http.Request, dst interface{}) error {
	mt, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		return httperror.New(http.StatusUnsupportedMediaType)
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body")
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			httperror.Message("Failed to unmarshal request body"),
			httperror.Cause(err),
		)
	}

	return nil
}

// JSONResponseBody sets content-type to application/json and marshals src
// as json to response body
// 500 Internal Server Error is returned if src could not be marshaled
func JSONResponseBody(w http.ResponseWriter, status int, src interface{}) error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	data, err := json.Marshal(src)
	if err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			httperror.Message("Failed to marshal response body"),
			httperror.Cause(err),
		)
	}

	w.WriteHeader(status)
	_, _ = w.Write(data)

	return nil
}

// DefaultRoutes sets up default responses for 404 Not Found and 405 MethodNotAllowed
func DefaultRoutes(router *httprouter.Router) {
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		data, err := json.Marshal(httperror.ErrorResponse{
			Message: http.StatusText(http.StatusNotFound),
		})
		if err != nil {
			log.Println(err)
		}

		_, _ = w.Write(data)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data, err := json.Marshal(httperror.ErrorResponse{
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
		if err != nil {
			log.Println(err)
		}

		_, _ = w.Write(data)
	})
}
