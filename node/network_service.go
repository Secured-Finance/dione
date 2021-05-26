package node

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/Secured-Finance/dione/node/wire"

	"github.com/Secured-Finance/dione/blockchain/pool"
)

const (
	MaxBlockCountForRetrieving = 500 // we do it just like in Bitcoin
)

type NetworkService struct {
	blockpool *pool.BlockPool
}

func NewNetworkService(bp *pool.BlockPool) *NetworkService {
	return &NetworkService{
		blockpool: bp,
	}
}

func (s *NetworkService) LastBlockHeight(ctx context.Context, arg interface{}, reply *wire.LastBlockHeightReply) {
	height, err := s.blockpool.GetLatestBlockHeight()
	if err != nil {
		reply.Error = err
		return
	}
	reply.Height = height
}

func (s *NetworkService) GetBlocks(ctx context.Context, arg wire.GetBlocksArg, reply *wire.GetBlocksReply) {
	if arg.From > arg.To {
		errText := "incorrect arguments: from > to"
		reply.Error = &errText
		return
	}
	if arg.To-arg.From > MaxBlockCountForRetrieving {
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
