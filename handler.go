package jsonrpc

import (
	"net/http"

	"github.com/intel-go/fastjson"
)

// Handler links a method of JSON-RPC request.
type Handler interface {
	ServeJSONRPC(c *Context, params *fastjson.RawMessage) (result interface{}, err *Error)
}

type HandlerChain []Handler

// ServeHTTP provides basic JSON-RPC handling.
func (mr *MethodRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rs, buf, batch, err := ParseRequest(r)
	if err != nil {
		err := SendResponse(w, []*Response{
			{
				Version: Version,
				Error:   err,
			},
		}, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := make([]*Response, len(rs))
	c := Context{r.Context(), buf}

	for i := range rs {
		if res := mr.InvokeMeddleware(c, rs[i]); res != nil {
			resp[i] = res
			continue
		}

		resp[i] = mr.InvokeMethod(c, rs[i])
	}

	if err := SendResponse(w, resp, batch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (mr *MethodRepository) InvokeMeddleware(c Context, r *Request) *Response {
	var res *Response

	for _, middleware := range mr.Middlewares {
		err := middleware(&c)
		if !c.IsNext() {
			if err != nil {
				res = NewResponse(r)
				res.Error = err
				res.Result = nil
			}
			break
		}
	}

	return res
}

// InvokeMethod invokes JSON-RPC method.
func (mr *MethodRepository) InvokeMethod(c Context, r *Request) *Response {
	var hs HandlerChain
	res := NewResponse(r)
	hs, res.Error = mr.TakeMethod(r)
	if res.Error != nil {
		return res
	}

	for _, h := range hs {
		ctx := WithRequestID(c, r.ID)
		res.Result, res.Error = h.ServeJSONRPC(&ctx, r.Params)
		if res.Error != nil {
			res.Result = nil
			break
		}
		if !ctx.IsNext() || res.Result != nil {
			break
		}
	}

	return res
}
