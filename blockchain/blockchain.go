package blockchain

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"

	"github.com/Secured-Finance/dione/blockchain/utils"

	types2 "github.com/Secured-Finance/dione/blockchain/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/ledgerwatch/lmdb-go/lmdb"
)

const (
	DefaultBlockDataPrefix   = "blockdata_"
	DefaultBlockHeaderPrefix = "header_"
	DefaultMetadataIndexName = "metadata"
	LatestBlockHeightKey     = "latest_block_height"
)

var (
	ErrBlockNotFound   = errors.New("block isn't found")
	ErrLatestHeightNil = errors.New("latest block height is nil")
)

type BlockChain struct {
	dbEnv         *lmdb.Env
	db            lmdb.DBI
	metadataIndex *utils.Index
	heightIndex   *utils.Index
}

func NewBlockChain(path string) (*BlockChain, error) {
	chain := &BlockChain{}

	// configure lmdb env
	env, err := lmdb.NewEnv()
	if err != nil {
		return nil, err
	}

	err = env.SetMaxDBs(1)
	if err != nil {
		return nil, err
	}
	err = env.SetMapSize(100 * 1024 * 1024 * 1024) // 100 GB
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	err = env.Open(path, 0, 0755)
	if err != nil {
		return nil, err
	}

	chain.dbEnv = env

	var dbi lmdb.DBI
	err = env.Update(func(txn *lmdb.Txn) error {
		dbi, err = txn.OpenDBI("blocks", lmdb.Create)
		return err
	})
	if err != nil {
		return nil, err
	}

	chain.db = dbi

	// create index instances
	metadataIndex := utils.NewIndex(DefaultMetadataIndexName, env, dbi)
	heightIndex := utils.NewIndex("height", env, dbi)
	chain.metadataIndex = metadataIndex
	chain.heightIndex = heightIndex

	return chain, nil
}

func (bp *BlockChain) setLatestBlockHeight(height uint64) error {
	return bp.metadataIndex.PutUint64([]byte(LatestBlockHeightKey), height)
}

func (bp *BlockChain) GetLatestBlockHeight() (uint64, error) {
	height, err := bp.metadataIndex.GetUint64([]byte(LatestBlockHeightKey))
	if err != nil {
		if err == utils.ErrIndexKeyNotFound {
			return 0, ErrLatestHeightNil
		}
		return 0, err
	}
	return height, nil
}

func (bp *BlockChain) StoreBlock(block *types2.Block) error {
	err := bp.dbEnv.Update(func(txn *lmdb.Txn) error {
		data, err := cbor.Marshal(block.Data)
		if err != nil {
			return err
		}
		headerData, err := cbor.Marshal(block.Header)
		if err != nil {
			return err
		}
		blockHash := hex.EncodeToString(block.Header.Hash)
		err = txn.Put(bp.db, []byte(DefaultBlockDataPrefix+blockHash), data, 0)
		if err != nil {
			return err
		}
		err = txn.Put(bp.db, []byte(DefaultBlockHeaderPrefix+blockHash), headerData, 0) // store header separately for easy fetching
		return err
	})
	if err != nil {
		return err
	}

	// update index "height -> block hash"
	heightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, block.Header.Height)
	err = bp.heightIndex.PutBytes(heightBytes, block.Header.Hash)
	if err != nil {
		return err
	}

	// update latest block height
	height, err := bp.GetLatestBlockHeight()
	if err != nil && err != ErrLatestHeightNil {
		return err
	}

	if err == ErrLatestHeightNil {
		if err = bp.setLatestBlockHeight(block.Header.Height); err != nil {
			return err
		}
	} else {
		if block.Header.Height > height {
			if err = bp.setLatestBlockHeight(block.Header.Height); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bp *BlockChain) HasBlock(blockHash []byte) (bool, error) {
	var blockExists bool
	err := bp.dbEnv.View(func(txn *lmdb.Txn) error {
		h := hex.EncodeToString(blockHash)
		_, err := txn.Get(bp.db, []byte(DefaultBlockHeaderPrefix+h)) // try to fetch block header
		if err != nil {
			if lmdb.IsNotFound(err) {
				blockExists = false
				return nil
			}
			return err
		}
		blockExists = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return blockExists, nil
}

func (bp *BlockChain) FetchBlockData(blockHash []byte) ([]*types2.Transaction, error) {
	var data []*types2.Transaction
	err := bp.dbEnv.View(func(txn *lmdb.Txn) error {
		h := hex.EncodeToString(blockHash)
		blockData, err := txn.Get(bp.db, []byte(DefaultBlockDataPrefix+h))
		if err != nil {
			if lmdb.IsNotFound(err) {
				return ErrBlockNotFound
			}
			return err
		}
		err = cbor.Unmarshal(blockData, &data)
		return err
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bp *BlockChain) FetchBlockHeader(blockHash []byte) (*types2.BlockHeader, error) {
	var blockHeader types2.BlockHeader
	err := bp.dbEnv.View(func(txn *lmdb.Txn) error {
		h := hex.EncodeToString(blockHash)
		data, err := txn.Get(bp.db, []byte(DefaultBlockHeaderPrefix+h))
		if err != nil {
			if lmdb.IsNotFound(err) {
				return ErrBlockNotFound
			}
			return err
		}
		err = cbor.Unmarshal(data, &blockHeader)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &blockHeader, nil
}

func (bp *BlockChain) FetchBlock(blockHash []byte) (*types2.Block, error) {
	var block types2.Block
	header, err := bp.FetchBlockHeader(blockHash)
	if err != nil {
		return nil, err
	}
	block.Header = header

	data, err := bp.FetchBlockData(blockHash)
	if err != nil {
		return nil, err
	}
	block.Data = data

	return &block, nil
}

func (bp *BlockChain) FetchBlockByHeight(height uint64) (*types2.Block, error) {
	var heightBytes = make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, height)
	blockHash, err := bp.heightIndex.GetBytes(heightBytes)
	if err != nil {
		if err == utils.ErrIndexKeyNotFound {
			return nil, ErrBlockNotFound
		}
	}
	block, err := bp.FetchBlock(blockHash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (bp *BlockChain) FetchBlockHeaderByHeight(height uint64) (*types2.BlockHeader, error) {
	var heightBytes = make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, height)
	blockHash, err := bp.heightIndex.GetBytes(heightBytes)
	if err != nil {
		if err == utils.ErrIndexKeyNotFound {
			return nil, ErrBlockNotFound
		}
	}
	blockHeader, err := bp.FetchBlockHeader(blockHash)
	if err != nil {
		return nil, err
	}
	return blockHeader, nil
}
