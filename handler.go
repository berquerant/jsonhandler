// Package jsonhandler provides functions to build http handler for JSON API.
package jsonhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

//go:generate go run github.com/berquerant/goconfig -type "SuccessStatusCode|int,MaxRequestBodyBytes|int64,ResponseContentCharset|string" -option -output config_generated.go -configOption Option

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
//
// Apply some options by pass Option to opt.
// Option is built by WithXXX functions.
//
// WithSuccessStatusCode changes the response http status code, default is 200 OK.
// WithMaxRequestBodyBytes limits the request body size, default is unlimited.
// WithResponseContentCharset changes the response content type charset, default is unspecified.
func Func[T, U any](handler Handler[T, U], errorHandler ErrorHandler, opt ...Option) http.HandlerFunc {
	config := NewConfigBuilder().
		SuccessStatusCode(http.StatusOK).
		MaxRequestBodyBytes(-1).
		ResponseContentCharset("").
		Build()
	config.Apply(opt...)

	iHandler := newInternalHandlerFunc(handler, config)

	return func(w http.ResponseWriter, r *http.Request) {
		if err := iHandler(w, r); err != nil && errorHandler != nil {
			errorHandler(w, r, err)
			return
		}
	}
}

var (
	errNotJSON = errors.New("not json")
)

func hasJSONHeader(header http.Header) bool {
	return strings.Contains(header.Get("Content-Type"), "application/json")
}

type internalHandlerFunc func(http.ResponseWriter, *http.Request) *Error

func newInternalHandlerFunc[T, U any](handler Handler[T, U], config *Config) internalHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		if !hasJSONHeader(r.Header) {
			return newError(errNotJSON, EnotJSONRequest)
		}
		t, err := readBody[T](w, r, config.MaxRequestBodyBytes.Get())
		if err != nil {
			return err
		}
		u, handleErr := handler(r.Context(), t)
		if handleErr != nil {
			return newError(handleErr, EhandlerError)
		}
		return writeBody[U](w, u, config.SuccessStatusCode.Get(), config.ResponseContentCharset.Get())
	}
}
