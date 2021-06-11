package types

import (
	"time"

	"github.com/Secured-Finance/dione/types"

	"github.com/libp2p/go-libp2p-core/crypto"

	"github.com/ethereum/go-ethereum/common"

	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"

	"github.com/libp2p/go-libp2p-core/peer"
)

type Block struct {
	Header *BlockHeader
	Data   []*Transaction
}

type BlockHeader struct {
	Timestamp     int64
	Height        uint64
	Hash          []byte
	LastHash      []byte
	LastHashProof *merkletree.Proof
	Proposer      peer.ID
	ProposerEth   common.Address
	Signature     []byte
	BeaconEntry   types.BeaconEntry
	ElectionProof *types.ElectionProof
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

func CreateBlock(lastBlockHeader *BlockHeader, txs []*Transaction, minerEth common.Address, privateKey crypto.PrivKey, eproof *types.ElectionProof) (*Block, error) {
	timestamp := time.Now().Unix()

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

	// sign the block hash
	s, err := privateKey.Sign(blockHash)
	if err != nil {
		return nil, err
	}

	lastHashProof, err := tree.GenerateProof(lastBlockHeader.Hash, 0)
	if err != nil {
		return nil, err
	}

	proposer, err := peer.IDFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Header: &BlockHeader{
			Timestamp:     timestamp,
			Height:        lastBlockHeader.Height + 1,
			Proposer:      proposer,
			ProposerEth:   minerEth,
			Signature:     s,
			Hash:          blockHash,
			LastHash:      lastBlockHeader.Hash,
			LastHashProof: lastHashProof,
			ElectionProof: eproof,
		},
		Data: txs,
	}

	return block, nil
}
