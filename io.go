package jsonhandler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// readBody tries to unmarshal body as T.
// this keeps request body.
func readBody[T any](r *http.Request) (t T, err *Error) {
	b, readErr := io.ReadAll(r.Body) // TODO: limit size
	if readErr != nil {
		err = newError(readErr, EreadRequestBody)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(b)) // keep body
	if jsonErr := json.Unmarshal(b, &t); jsonErr != nil {
		err = newError(jsonErr, EunmarshalRequestBody)
	}
	return
}

// writeBody tries to write the json encoding of t.
func writeBody[T any](w http.ResponseWriter, t T) *Error {
	b, err := json.Marshal(t)
	if err != nil {
		return newError(err, EmarshalResponse)
	}
	w.Header().Add("Content-Type", "application/json") // TODO: specify charset
	if _, err := w.Write(b); err != nil {
		return newError(err, EwriteResponseBody)
	}
	return nil
}
