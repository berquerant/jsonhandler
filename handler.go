// Package jsonhandler provides functions to build http handler for JSON API.
package jsonhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type (
	// ErrorHandler handles an error on http request.
	ErrorHandler func(http.ResponseWriter, *http.Request, *Error)
	// Handler responds to a request of type T.
	Handler[T, U any] func(context.Context, T) (U, error)
)

// Func builds a http handler for JSON API.
//
// Func calls handler with the context of the request and the parsed request body,
// also calls errorHandler when an error occurred and errorHandler is not nil.
func Func[T, U any](handler Handler[T, U], errorHandler ErrorHandler) http.HandlerFunc {
	iHandler := newInternalHandlerFunc(handler)

	return func(w http.ResponseWriter, r *http.Request) {
		if err := iHandler(w, r); err != nil && errorHandler != nil {
			errorHandler(w, r, err)
		}
		// TODO: change status code
	}
}

var (
	errNotJSON = errors.New("not json")
)

func hasJSONHeader(header http.Header) bool {
	return strings.Contains(header.Get("Content-Type"), "application/json")
}

type internalHandlerFunc func(http.ResponseWriter, *http.Request) *Error

func newInternalHandlerFunc[T, U any](handler Handler[T, U]) internalHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		if !hasJSONHeader(r.Header) {
			return newError(errNotJSON, EnotJSONRequest)
		}
		t, err := readBody[T](r)
		if err != nil {
			return err
		}
		u, handleErr := handler(r.Context(), t)
		if handleErr != nil {
			return newError(handleErr, EhandlerError)
		}
		return writeBody[U](w, u)
	}
}
