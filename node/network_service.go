package node

import (
	"context"
	"errors"
	"fmt"

	gorpc "github.com/libp2p/go-libp2p-gorpc"

	"github.com/sirupsen/logrus"

	"github.com/Secured-Finance/dione/consensus/policy"
	"github.com/Secured-Finance/dione/node/wire"

	"github.com/Secured-Finance/dione/blockchain/pool"
)

type NetworkService struct {
	blockpool *pool.BlockPool
	mempool   *pool.Mempool
	rpcClient *gorpc.Client
}

func NewNetworkService(bp *pool.BlockPool) *NetworkService {
	return &NetworkService{
		blockpool: bp,
	}
}

func (s *NetworkService) LastBlockHeight(ctx context.Context, arg struct{}, reply *wire.LastBlockHeightReply) {
	height, err := s.blockpool.GetLatestBlockHeight()
	if err != nil {
		reply.Error = err
		return
	}
	reply.Height = height
}

func (s *NetworkService) GetRangeOfBlocks(ctx context.Context, arg wire.GetRangeOfBlocksArg, reply *wire.GetRangeOfBlocksReply) {
	if arg.From > arg.To {
		errText := "incorrect arguments: from > to"
		reply.Error = &errText
		return
	}
	if arg.To-arg.From > policy.MaxBlockCountForRetrieving {
		errText := "incorrect arguments: count of block for retrieving is exceeded the limit"
		reply.Error = &errText
		return
	}
	for i := arg.From; i <= arg.To; i++ {
		block, err := s.blockpool.FetchBlockByHeight(i)
		if err != nil {
			logrus.Warnf("failed to retrieve block from blockpool with height %d", i)
			reply.FailedBlockHeights = append(reply.FailedBlockHeights, i)
			continue
		}
		reply.Blocks = append(reply.Blocks, *block)
	}
}

func (s *NetworkService) Mempool(ctx context.Context, arg struct{}, reply *wire.InvMessage) {
	txs := s.mempool.GetAllTransactions()

	// extract hashes of txs
	for _, v := range txs {
		reply.Inventory = append(reply.Inventory, wire.InvItem{
			Type: wire.TxInvType,
			Hash: v.Hash,
		})
	}
}

func (s *NetworkService) GetMempoolTxs(ctx context.Context, arg wire.GetMempoolTxsArg, reply *wire.GetMempoolTxsReply) {
	if len(arg.Items) > MaxTransactionCountForRetrieving {
		pid, _ := gorpc.GetRequestSender(ctx)
		logrus.Warnf("Max tx count limit exceeded for GetMempoolTxs request of node %s", pid)
		reply.Error = fmt.Errorf("max tx count limit exceeded")
		return
	}

	for _, v := range arg.Items {
		tx, err := s.mempool.GetTransaction(v)
		if err != nil {
			if errors.Is(err, pool.ErrTxNotFound) {
				reply.NotFoundTxs = append(reply.NotFoundTxs, v)
			} else {
				reply.Error = err
				return
			}
		}
		reply.Transactions = append(reply.Transactions, *tx)
	}
}
