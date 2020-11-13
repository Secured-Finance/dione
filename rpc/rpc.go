package rpc

import "net/http"

type RequestBody struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

func NewRequestBody(method string) *RequestBody {
	return &RequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  []string{},
		ID:      0,
	}
}

type Client interface {
	HandleRequest(r *http.Request, data []byte) (*http.Response, error)
}
