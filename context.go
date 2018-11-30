package jsonrpc

import (
	"context"

	"github.com/intel-go/fastjson"
)

const NEXT = "NEXT"

type requestIDKey struct{}

type Context struct {
	ctx  context.Context
	body []byte
}

// WithRequestID adds request id to context.
func WithRequestID(c Context, id *fastjson.RawMessage) Context {
	return Context{context.WithValue(c.ctx, requestIDKey{}, id), c.body}
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Body() []byte {
	return c.body
}

// RequestID takes request id from context.
func (c *Context) RequestID() *fastjson.RawMessage {
	return c.ctx.Value(requestIDKey{}).(*fastjson.RawMessage)
}

func (c *Context) Next() {
	c.ctx = context.WithValue(c.ctx, NEXT, true)
}

func (c *Context) Abort() {
	c.ctx = context.WithValue(c.ctx, NEXT, false)
}

//default is true
func (c *Context) IsNext() bool {
	if c.ctx.Value(NEXT) == nil {
		return true
	}
	return c.ctx.Value(NEXT).(bool)
}
