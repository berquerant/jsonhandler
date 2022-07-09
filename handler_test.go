package jsonhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berquerant/jsonhandler"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type funcTestcase[T, U any] struct {
	title    string
	handler  jsonhandler.Handler[T, U]
	header   map[string]string
	body     string
	request  T
	response U
	ok       bool
	errKind  jsonhandler.ErrorKind
}

func newFuncTestcase[T, U any](
	title string,
	handler jsonhandler.Handler[T, U],
	header map[string]string,
	body string,
	request T,
	response U,
	ok bool,
	errKind jsonhandler.ErrorKind,
) func(*testing.T) {
	x := &funcTestcase[T, U]{
		title:    title,
		handler:  handler,
		header:   header,
		body:     body,
		request:  request,
		response: response,
		ok:       ok,
		errKind:  errKind,
	}
	return func(t *testing.T) {
		t.Run(x.title, x.test)
	}
}

func (s *funcTestcase[T, U]) test(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com",
		bytes.NewBufferString(s.body))
	for k, v := range s.header {
		req.Header.Add(k, v)
	}
	rec := httptest.NewRecorder()

	var respErr *jsonhandler.Error
	onError := func(_ http.ResponseWriter, _ *http.Request, err *jsonhandler.Error) {
		respErr = err
	}

	var gotResponse U
	handler := func(ctx context.Context, r T) (U, error) {
		assert.Equal(t, cmp.Diff(s.request, r), "", "request")
		u, err := s.handler(ctx, r)
		gotResponse = u
		return u, err
	}

	jsonhandler.Func(handler, onError)(rec, req)

	if !s.ok {
		if !assert.NotNil(t, respErr, "response error is not nil") {
			return
		}
		assert.Equal(t, s.errKind, respErr.Kind, "repsonse error kind")
		return
	}

	assert.Nil(t, respErr, "response error is nil")
	assert.Equal(t, cmp.Diff(s.response, gotResponse), "", "response")

	resp := rec.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "http status")
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json", "content-type")

	var writtenResponse U
	respBody, err := io.ReadAll(resp.Body)
	if !assert.Nil(t, err) {
		return
	}
	assert.Nil(t, json.Unmarshal(respBody, &writtenResponse), "writte response unmarshal")
	assert.Equal(t, cmp.Diff(s.response, writtenResponse), "", "written response")
}

func newHandler[T, U any](u U, err error) jsonhandler.Handler[T, U] {
	return func(_ context.Context, t T) (U, error) {
		return u, err
	}
}

type (
	RequestT struct {
		Message string `json:"message"`
	}
	ResponseT struct {
		Status string `json:"status"`
	}

	DictT = map[string]any
)

func TestFunc(t *testing.T) {
	const (
		requestTJSON        = `{"message":"hello"}`
		responseTJSON       = `{"status":"ok"}`
		invalidRequestTJSON = `not json`
	)

	var (
		request = RequestT{
			Message: "hello",
		}
		response = ResponseT{
			Status: "ok",
		}
		normalHandler = newHandler[RequestT, ResponseT](response, nil)
		errorHandler  = newHandler[RequestT, ResponseT](response, errors.New("handler error"))
		normalHeader  = map[string]string{
			"Content-Type": "application/json",
		}
		newTestcase = newFuncTestcase[RequestT, ResponseT]

		dictRequest = map[string]any{
			"message": "hello",
		}
		dictResponse = map[string]any{
			"status": "end",
		}
	)

	for _, tc := range []func(*testing.T){
		newTestcase(
			"reject request without application/json in header",
			normalHandler,
			nil,
			requestTJSON,
			request,
			response,
			false,
			jsonhandler.EnotJSONRequest,
		),
		newTestcase(
			"reject malformed request body",
			normalHandler,
			normalHeader,
			invalidRequestTJSON,
			request,
			response,
			false,
			jsonhandler.EunmarshalRequestBody,
		),
		newFuncTestcase[RequestT, func()](
			"reject the response that cannot be jsonified",
			newHandler[RequestT, func()](func() {}, nil),
			normalHeader,
			requestTJSON,
			request,
			func() {},
			false,
			jsonhandler.EmarshalResponse,
		),
		newTestcase(
			"detect the handler error",
			errorHandler,
			normalHeader,
			requestTJSON,
			request,
			response,
			false,
			jsonhandler.EhandlerError,
		),
		newTestcase(
			"call handler with struct parameter successfully",
			normalHandler,
			normalHeader,
			requestTJSON,
			request,
			response,
			true,
			jsonhandler.Eunknown,
		),
		newFuncTestcase[*RequestT, *ResponseT](
			"call handler with pointer parameter successfully",
			newHandler[*RequestT, *ResponseT](&response, nil),
			normalHeader,
			requestTJSON,
			&request,
			&response,
			true,
			jsonhandler.Eunknown,
		),
		newFuncTestcase[DictT, DictT](
			"call handler with map parameter successfully",
			newHandler[DictT, DictT](dictResponse, nil),
			normalHeader,
			requestTJSON,
			dictRequest,
			dictResponse,
			true,
			jsonhandler.Eunknown,
		),
	} {
		tc(t)
	}
}
