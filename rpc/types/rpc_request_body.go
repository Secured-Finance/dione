package types

type RPCRequestBody struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func NewRPCRequestBody(method string) *RPCRequestBody {
	var i []interface{}
	return &RPCRequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  i,
		ID:      0,
	}
}
