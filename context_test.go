package jsonrpc

import (
	"context"
	"testing"

	"github.com/intel-go/fastjson"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {

	c := Context{context.Background(), nil}
	id := fastjson.RawMessage("1")
	c = WithRequestID(c, &id)
	var pick *fastjson.RawMessage
	require.NotPanics(t, func() {
		pick = c.RequestID()
	})
	require.Equal(t, &id, pick)
}

func TestNext(t *testing.T) {
	c := Context{context.Background(), nil}
	var pick bool
	require.NotPanics(t, func() {
		pick = c.IsNext()
	})
	require.Equal(t, true, pick)
	c.Next()
	require.NotPanics(t, func() {
		pick = c.IsNext()
	})
	require.Equal(t, true, pick)
	c.Abort()
	require.NotPanics(t, func() {
		pick = c.IsNext()
	})
	require.Equal(t, false, pick)
}
