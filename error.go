package jsonhandler

import (
	"encoding/json"
	"fmt"
)

type ErrorKind int

//go:generate go run golang.org/x/tools/cmd/stringer -type ErrorKind -output errorkind_generated.go

const (
	Eunknown ErrorKind = iota
	// EnotJSONRequest indicates that the request body is not json.
	EnotJSONRequest
	// EreadRequestBody indicates that the server failed to read the request body.
	EreadRequestBody
	// EunmarshalRequestBody indicates that the request body cannot be parsed.
	EunmarshalRequestBody
	// EmarshalResponse indicates that the server failed to encode the response into json.
	EmarshalResponse
	// EwriteResponseBody indicates that the server failed to write the response body.
	EwriteResponseBody
	// EhandlerError indicates that the Handler returned an error.
	EhandlerError
	// EtooLargeRequestBody indicates that the request body is too large.
	EtooLargeRequestBody
)

// IsHandlerError returns this is EhandlerError or not.
func (ek ErrorKind) IsHandlerError() bool { return ek == EhandlerError }

type Error struct {
	// Err is the original error.
	Err  error
	Kind ErrorKind
}

func (err *Error) Error() string { return fmt.Sprintf("%s %v", err.Kind, err.Err) }

func (err *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"err":  err.Err.Error(),
		"kind": err.Kind.String(),
	})
}

func newError(err error, kind ErrorKind) *Error {
	return &Error{
		Err:  err,
		Kind: kind,
	}
}
