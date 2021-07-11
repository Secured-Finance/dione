package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/Secured-Finance/dione/blockchain"

	gorpc "github.com/libp2p/go-libp2p-gorpc"

	"github.com/sirupsen/logrus"

	"github.com/Secured-Finance/dione/consensus/policy"
	"github.com/Secured-Finance/dione/node/wire"

	"github.com/Secured-Finance/dione/blockchain/pool"
)

type NetworkService struct {
	blockchain *blockchain.BlockChain
	mempool    *pool.Mempool
	rpcClient  *gorpc.Client
}

func NewNetworkService(bc *blockchain.BlockChain, mp *pool.Mempool) *NetworkService {
	return &NetworkService{
		blockchain: bc,
		mempool:    mp,
	}
}

func (s *NetworkService) LastBlockHeight(ctx context.Context, arg struct{}, reply *wire.LastBlockHeightReply) error {
	height, err := s.blockchain.GetLatestBlockHeight()
	if err != nil {
		return err
	}
	reply.Height = height
	return nil
}

func (s *NetworkService) GetRangeOfBlocks(ctx context.Context, arg wire.GetRangeOfBlocksArg, reply *wire.GetRangeOfBlocksReply) error {
	if arg.From > arg.To {
		return fmt.Errorf("incorrect arguments: from > to")
	}
	if arg.To-arg.From > policy.MaxBlockCountForRetrieving {
		return fmt.Errorf("incorrect arguments: count of block for retrieving is exceeded the limit")
	}
	for i := arg.From; i <= arg.To; i++ {
		block, err := s.blockchain.FetchBlockByHeight(i)
		if err != nil {
			logrus.Warnf("failed to retrieve block from blockpool with height %d", i)
			reply.FailedBlockHeights = append(reply.FailedBlockHeights, i)
			continue
		}
		reply.Blocks = append(reply.Blocks, *block)
	}
	return nil
}

func (s *NetworkService) Mempool(ctx context.Context, arg struct{}, reply *wire.InvMessage) error {
	txs := s.mempool.GetAllTransactions()

	// extract hashes of txs
	for _, v := range txs {
		reply.Inventory = append(reply.Inventory, wire.InvItem{
			Type: wire.TxInvType,
			Hash: v.Hash,
		})
	}

	return nil
}

func (s *NetworkService) GetMempoolTxs(ctx context.Context, arg wire.GetMempoolTxsArg, reply *wire.GetMempoolTxsReply) error {
	if len(arg.Items) > policy.MaxTransactionCountForRetrieving {
		pid, _ := gorpc.GetRequestSender(ctx)
		logrus.Warnf("Max tx count limit exceeded for GetMempoolTxs request of node %s", pid)
		return fmt.Errorf("max tx count limit exceeded")
	}

	for _, v := range arg.Items {
		tx, err := s.mempool.GetTransaction(v)
		if err != nil {
			if errors.Is(err, pool.ErrTxNotFound) {
				reply.NotFoundTxs = append(reply.NotFoundTxs, v)
			} else {
				return err
			}
		}
		reply.Transactions = append(reply.Transactions, *tx)
	}

	return nil
}
