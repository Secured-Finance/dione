package types

type RPCResponseBody struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  []byte `json:"result"`
	Error   Error  `json:"error"`
	ID      int64  `json:"id"`
}
