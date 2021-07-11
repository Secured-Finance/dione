package wire

import "github.com/Secured-Finance/dione/blockchain/types"

type GetRangeOfBlocksArg struct {
	From uint64
	To   uint64
}

type GetRangeOfBlocksReply struct {
	Blocks             []types.Block
	FailedBlockHeights []uint64 // list of block heights the node was unable to retrieve
}
