package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/intel-go/fastjson"
)

type (
	EchoHandler struct{}
	EchoParams  struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}

	MiddlewareHandler struct{}

	PositionalHandler struct{}
	PositionalParams  []int
	PositionalResult  struct {
		Message []int `json:"message"`
	}
)

func (h EchoHandler) ServeJSONRPC(c *Context, params *fastjson.RawMessage) (interface{}, *Error) {

	var p EchoParams
	if err := Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func (h MiddlewareHandler) ServeJSONRPC(c *Context, params *fastjson.RawMessage) (interface{}, *Error) {

	var p PositionalParams
	if err := Unmarshal(params, &p); err != nil {
		return nil, err
	}
	log.Printf("%+v", p)

	return nil, nil
}

func (h PositionalHandler) ServeJSONRPC(c *Context, params *fastjson.RawMessage) (interface{}, *Error) {

	var p PositionalParams
	if err := Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return PositionalResult{
		Message: p,
	}, nil
}

func ExampleEchoHandler_ServeJSONRPC() {

	mr := NewMethodRepository()

	mr.Use(func(context *Context) (err *Error) {
		req := map[string]interface{}{}
		if err := json.Unmarshal(context.Body(), &req); err != nil {
			return ErrParse()
		}
		log.Printf("middleware %s", req)
		context.Abort()
		return nil
	})

	if err := mr.RegisterMethod("Main.Echo", EchoParams{}, EchoResult{}, EchoHandler{}); err != nil {
		log.Fatalln(err)
	}

	if err := mr.RegisterMethod("Main.Positional", PositionalParams{}, PositionalResult{},
		MiddlewareHandler{}, PositionalHandler{}); err != nil {
		log.Fatalln(err)
	}

	http.Handle("/jrpc", mr)
	http.HandleFunc("/jrpc/debug", mr.ServeDebug)

	srv := httptest.NewServer(http.DefaultServeMux)
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/jrpc", "application/json", bytes.NewBufferString(`{
		"jsonrpc": "2.0",
		"method": "Main.Echo",
		"params": {
			"name": "John Doe"
		},
		"id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
	}`))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if _, err = io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatalln(err)
	}

	resp, err = http.Post(srv.URL+"/jrpc", "application/json", bytes.NewBufferString(`{
		"jsonrpc": "2.0",
		"method": "Main.Positional",
		"params": [3, 1, 1, 3, 5, 3],
		"id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
	}`))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// {"jsonrpc":"2.0","result":{"message":"Hello, John Doe"},"id":"243a718a-2ebb-4e32-8cc8-210c39e8a14b"}
	// {"jsonrpc":"2.0","result":{"message":[3,1,1,3,5,3]},"id":"243a718a-2ebb-4e32-8cc8-210c39e8a14b"}
}
