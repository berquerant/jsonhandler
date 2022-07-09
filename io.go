package jsonhandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newMaxBytesReader(w http.ResponseWriter, r io.ReadCloser, maxBytes int64) io.Reader {
	if maxBytes < 0 {
		return r
	}
	return http.MaxBytesReader(w, r, maxBytes)
}

// readBody tries to unmarshal body as T.
// this keeps request body.
// maxBodyBytes limits body size, negative value means unlimited.
func readBody[T any](w http.ResponseWriter, r *http.Request, maxBodyBytes int64) (t T, err *Error) {
	b, readErr := io.ReadAll(newMaxBytesReader(w, r.Body, maxBodyBytes))
	if readErr != nil {
		if readErr.Error() == "http: request body too large" {
			err = newError(readErr, EtooLargeRequestBody)
		} else {
			err = newError(readErr, EreadRequestBody)
		}
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(b)) // keep body
	if jsonErr := json.Unmarshal(b, &t); jsonErr != nil {
		err = newError(jsonErr, EunmarshalRequestBody)
	}
	return
}

func newResponseContentType(charset string) string {
	const base = "application/json"
	if charset != "" {
		return fmt.Sprintf("%s; charset=%s", base, charset)
	}
	return base
}

// writeBody tries to write the json encoding of t.
// charset is for response Content-Type header.
func writeBody[T any](w http.ResponseWriter, t T, statusCode int, charset string) *Error {
	b, err := json.Marshal(t)
	if err != nil {
		return newError(err, EmarshalResponse)
	}
	w.Header().Add("Content-Type", newResponseContentType(charset))
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		return newError(err, EwriteResponseBody)
	}
	return nil
}
