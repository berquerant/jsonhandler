# jsonhandler

Package jsonhandler provides functions to build http handler for JSON API.


## Examples

``` go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/berquerant/jsonhandler"
)

type (
	Request struct {
		Message string
	}
	Response struct {
		Result string
	}
)

func handler(_ context.Context, req *Request) (*Response, error) {
	return &Response{
		Result: fmt.Sprintf("I got `%s`", req.Message),
	}, nil
}

func main() {
	http.HandleFunc("/", jsonhandler.Func(handler, nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

``` shell
❯ curl -s localhost:8080/ -H 'content-type: application/json' -d '{"Message": "Hello!"}'
{"Result":"I got `Hello!`"}
```
