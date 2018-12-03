package jsonrpc

import (
	"context"

	"github.com/intel-go/fastjson"
)

type middlewareNextKey struct{}

type requestIDKey struct{}

//Context body is http post data
type Context struct {
	ctx  context.Context
	body []byte
}

// WithRequestID adds request id to context.
func WithRequestID(c Context, id *fastjson.RawMessage) Context {
	return Context{context.WithValue(c.ctx, requestIDKey{}, id), c.body}
}

//Context get context
func (c *Context) Context() context.Context {
	return c.ctx
}

//Body get body
func (c *Context) Body() []byte {
	return c.body
}

// RequestID takes request id from context.
func (c *Context) RequestID() *fastjson.RawMessage {
	return c.ctx.Value(requestIDKey{}).(*fastjson.RawMessage)
}

//Next middleware invokes, JSON-RPC continue
func (c *Context) Next() {
	c.ctx = context.WithValue(c.ctx, middlewareNextKey{}, true)
}

//Abort middleware invokes, JSON-RPC stop
func (c *Context) Abort() {
	c.ctx = context.WithValue(c.ctx, middlewareNextKey{}, false)
}

//IsNext  can the JSON-RPC continue,  default is true
func (c *Context) IsNext() bool {
	if c.ctx.Value(middlewareNextKey{}) == nil {
		return true
	}
	return c.ctx.Value(middlewareNextKey{}).(bool)
}
