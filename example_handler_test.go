package jsonhandler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/berquerant/jsonhandler"
)

type (
	intCounter struct {
		mux sync.RWMutex
		val int
	}

	CounterMutateRequest struct {
		Op    string `json:"op"`
		Value int    `json:"value"`
	}
	CounterMutateResponse struct {
		Value int `json:"value"`
	}
)

func (c *intCounter) get() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.val
}

func (c *intCounter) delta(n int) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.val += n
	return c.val
}

func (c *intCounter) fix(n int) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.val = n
	return c.val
}

func (c *intCounter) handler(_ context.Context, r *CounterMutateRequest) (*CounterMutateResponse, error) {
	switch r.Op {
	case "delta":
		return &CounterMutateResponse{
			Value: c.delta(r.Value),
		}, nil
	case "fix":
		return &CounterMutateResponse{
			Value: c.fix(r.Value),
		}, nil
	case "get":
		return &CounterMutateResponse{
			Value: c.get(),
		}, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", r.Op)
	}
}

func onError(w http.ResponseWriter, r *http.Request, err *jsonhandler.Error) {
	switch err.Kind {
	case jsonhandler.EnotJSONRequest, jsonhandler.EunmarshalRequestBody, jsonhandler.EhandlerError:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	b, _ := io.ReadAll(r.Body)
	m, _ := json.Marshal(map[string]any{
		"body": string(b),
		"err":  err,
	})
	_, _ = w.Write(m)
}

func ExampleFunc() {
	counter := new(intCounter)
	http.HandleFunc("/c", jsonhandler.Func(counter.handler, onError))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
