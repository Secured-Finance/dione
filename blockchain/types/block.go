package types

import (
	"time"

	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"

	"github.com/Secured-Finance/dione/wallet"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Block struct {
	Header *BlockHeader
	Data   []*Transaction
}

type BlockHeader struct {
	Timestamp int64
	Height    uint64
	Hash      []byte
	LastHash  []byte
	Proposer  peer.ID
	Signature []byte
}

func GenesisBlock() *Block {
	return &Block{
		Header: &BlockHeader{
			Timestamp: 1620845070,
			Height:    0,
			Hash:      []byte("DIMICANDUM"),
		},
		Data: []*Transaction{},
	}
}

func CreateBlock(lastBlockHeader *BlockHeader, txs []*Transaction, wallet *wallet.LocalWallet) (*Block, error) {
	timestamp := time.Now().Unix()
	proposer, err := wallet.GetDefault()
	if err != nil {
		return nil, err
	}

	// extract hashes from transactions
	var txHashes [][]byte
	for _, tx := range txs {
		txHashes = append(txHashes, tx.Hash)
	}
	txHashes = append(txHashes, lastBlockHeader.Hash)

	tree, err := merkletree.NewUsing(txHashes, keccak256.New(), false)
	if err != nil {
		return nil, err
	}

	// fetch merkle tree root hash (block hash)
	blockHash := tree.Root()

	// sign this block hash
	s, err := wallet.Sign(proposer, blockHash)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Header: &BlockHeader{
			Timestamp: timestamp,
			Height:    lastBlockHeader.Height + 1,
			Proposer:  proposer,
			Signature: s.Data,
			Hash:      blockHash,
			LastHash:  lastBlockHeader.Hash,
		},
		Data: txs,
	}

	return block, nil
}
