package types

import (
	types2 "github.com/Secured-Finance/dione/blockchain/types"
)

type PrePrepareMessage struct {
	Block *types2.Block
}

type PrepareMessage struct {
	Blockhash []byte
	Signature []byte
}

type CommitMessage PrepareMessage
