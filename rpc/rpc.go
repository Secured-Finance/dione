package rpc

import "net/http"

type RequestBody struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func NewRequestBody(method string) *RequestBody {
	var i []interface{}
	return &RequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  i,
		ID:      0,
	}
}

type Client interface {
	HandleRequest(r *http.Request, data []byte) (*http.Response, error)
}
