package rpc

type RPCClient interface {
	GetTransaction(txHash string) ([]byte, error)
}
