package pool

import (
	"encoding/hex"

	"github.com/Secured-Finance/dione/blockchain/types"
	"github.com/Secured-Finance/dione/cache"
)

// BlockPool is pool for blocks that isn't not validated or committed yet
type BlockPool struct {
	mempool        *Mempool
	knownBlocks    cache.Cache
	acceptedBlocks cache.Cache
}

func NewBlockPool(mp *Mempool) (*BlockPool, error) {
	bp := &BlockPool{
		acceptedBlocks: cache.NewInMemoryCache(), // here we need to use separate cache
		knownBlocks:    cache.NewInMemoryCache(),
		mempool:        mp,
	}

	return bp, nil
}

func (bp *BlockPool) AddBlock(block *types.Block) error {
	return bp.knownBlocks.Store(hex.EncodeToString(block.Header.Hash), block)
}

func (bp *BlockPool) GetBlock(blockhash []byte) (*types.Block, error) {
	var block *types.Block
	return block, bp.knownBlocks.Get(hex.EncodeToString(blockhash), &block)
}

// PruneBlocks cleans known blocks list. It is called when new consensus round starts.
func (bp *BlockPool) PruneBlocks() {
	for k := range bp.knownBlocks.Items() {
		bp.knownBlocks.Delete(k)
	}
}

func (bp *BlockPool) AddAcceptedBlock(block *types.Block) error {
	return bp.acceptedBlocks.Store(hex.EncodeToString(block.Header.Hash), block)
}

func (bp *BlockPool) GetAllAcceptedBlocks() []*types.Block {
	var blocks []*types.Block
	for _, v := range bp.acceptedBlocks.Items() {
		blocks = append(blocks, v.(*types.Block))
	}
	return blocks
}

// PruneAcceptedBlocks cleans accepted blocks list. It is called when new consensus round starts.
func (bp *BlockPool) PruneAcceptedBlocks() {
	for k, v := range bp.acceptedBlocks.Items() {
		block := v.(*types.Block)
		for _, v := range block.Data {
			v.MerkleProof = nil
			bp.mempool.StoreTx(v) // return transactions back to mempool
		}
		bp.acceptedBlocks.Delete(k)
	}
}
