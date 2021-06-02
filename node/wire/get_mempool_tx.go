package wire

import "github.com/Secured-Finance/dione/blockchain/types"

type GetMempoolTxsArg struct {
	Items [][]byte
}

type GetMempoolTxsReply struct {
	Transactions []types.Transaction
	NotFoundTxs  [][]byte
	Error        error
}
