package consensus

import (
	"errors"
	"math/big"
	"sync"

	"github.com/fxamacker/cbor/v2"

	"github.com/Secured-Finance/dione/cache"

	"github.com/asaskevich/EventBus"

	"github.com/Secured-Finance/dione/blockchain"

	"github.com/drand/drand/client"

	"github.com/Arceliar/phony"

	types3 "github.com/Secured-Finance/dione/blockchain/types"

	"github.com/libp2p/go-libp2p-core/crypto"

	"github.com/Secured-Finance/dione/blockchain/pool"

	"github.com/Secured-Finance/dione/consensus/types"
	types2 "github.com/Secured-Finance/dione/types"

	"github.com/Secured-Finance/dione/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/Secured-Finance/dione/pubsub"
)

var (
	ErrNoAcceptedBlocks = errors.New("there is no accepted blocks")
)

type StateStatus uint8

const (
	StateStatusUnknown = iota

	StateStatusPrePrepared
	StateStatusPrepared
	StateStatusCommited
)

type PBFTConsensusManager struct {
	phony.Inbox
	bus            EventBus.Bus
	psb            *pubsub.PubSubRouter
	minApprovals   int // FIXME
	privKey        crypto.PrivKey
	msgLog         *ConsensusMessageLog
	validator      *ConsensusValidator
	ethereumClient *ethclient.EthereumClient
	miner          *Miner
	blockPool      *pool.BlockPool
	blockchain     blockchain.BlockChain
	state          *State
}

type State struct {
	mutex       sync.Mutex
	drandRound  uint64
	randomness  []byte
	blockHeight uint64
	status      StateStatus
	ready       chan bool
}

func NewPBFTConsensusManager(
	bus EventBus.Bus,
	psb *pubsub.PubSubRouter,
	minApprovals int,
	privKey crypto.PrivKey,
	ethereumClient *ethclient.EthereumClient,
	miner *Miner,
	bc *blockchain.BlockChain,
	bp *pool.BlockPool,
) *PBFTConsensusManager {
	pcm := &PBFTConsensusManager{}
	pcm.psb = psb
	pcm.miner = miner
	pcm.validator = NewConsensusValidator(miner, bc)
	pcm.msgLog = NewConsensusMessageLog()
	pcm.minApprovals = minApprovals
	pcm.privKey = privKey
	pcm.ethereumClient = ethereumClient
	pcm.state = &State{
		ready:  make(chan bool, 1),
		status: StateStatusUnknown,
	}
	pcm.bus = bus
	pcm.blockPool = bp
	pcm.psb.Hook(pubsub.PrePrepareMessageType, pcm.handlePrePrepare, types.PrePrepareMessage{})
	pcm.psb.Hook(pubsub.PrepareMessageType, pcm.handlePrepare, types.PrepareMessage{})
	pcm.psb.Hook(pubsub.CommitMessageType, pcm.handleCommit, types.CommitMessage{})
	bus.SubscribeOnce("sync:initialSyncCompleted", func() {
		pcm.state.ready <- true
	})
	return pcm
}

func (pcm *PBFTConsensusManager) propose(blk *types3.Block) error {
	pcm.state.mutex.Lock()
	defer pcm.state.mutex.Unlock()
	prePrepareMsg, err := NewMessage(types.ConsensusMessage{Block: blk}, types.ConsensusMessageTypePrePrepare, pcm.privKey)
	if err != nil {
		return err
	}
	pcm.psb.BroadcastToServiceTopic(prePrepareMsg)
	pcm.blockPool.AddBlock(blk)
	pcm.state.status = StateStatusPrePrepared
	return nil
}

func (pcm *PBFTConsensusManager) handlePrePrepare(message *pubsub.GenericMessage) {
	pcm.state.mutex.Lock()
	defer pcm.state.mutex.Unlock()
	prePrepare, ok := message.Payload.(types.PrePrepareMessage)
	if !ok {
		logrus.Warn("failed to convert payload to PrePrepare message")
		return
	}

	if prePrepare.Block.Header.Proposer == pcm.miner.address {
		return
	}

	cmsg := types.ConsensusMessage{
		Type:  types.ConsensusMessageTypePrePrepare,
		From:  message.From,
		Block: prePrepare.Block,
	}

	<-pcm.state.ready

	if pcm.msgLog.Exists(cmsg) {
		logrus.Debugf("received existing pre_prepare msg, dropping...")
		return
	}
	if !pcm.validator.Valid(cmsg, map[string]interface{}{"randomness": pcm.state.randomness}) {
		logrus.Warn("received invalid pre_prepare msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(cmsg)
	pcm.blockPool.AddBlock(cmsg.Block)

	prepareMsg, err := NewMessage(cmsg, types.ConsensusMessageTypePrepare, pcm.privKey)
	if err != nil {
		logrus.Errorf("failed to create prepare message: %v", err)
		return
	}

	pcm.psb.BroadcastToServiceTopic(prepareMsg)
	pcm.state.status = StateStatusPrePrepared
}

func (pcm *PBFTConsensusManager) handlePrepare(message *pubsub.GenericMessage) {
	pcm.state.mutex.Lock()
	defer pcm.state.mutex.Unlock()
	prepare, ok := message.Payload.(types.PrepareMessage)
	if !ok {
		logrus.Warn("failed to convert payload to Prepare message")
		return
	}

	cmsg := types.ConsensusMessage{
		Type:      types.ConsensusMessageTypePrepare,
		From:      message.From,
		Blockhash: prepare.Blockhash,
		Signature: prepare.Signature,
	}

	if _, err := pcm.blockPool.GetBlock(cmsg.Blockhash); errors.Is(err, cache.ErrNotFound) {
		logrus.Debugf("received unknown block")
		return
	}

	if pcm.msgLog.Exists(cmsg) {
		logrus.Debugf("received existing prepare msg, dropping...")
		return
	}

	if !pcm.validator.Valid(cmsg, nil) {
		logrus.Warn("received invalid prepare msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(cmsg)

	if len(pcm.msgLog.Get(types.ConsensusMessageTypePrepare, cmsg.Blockhash)) >= pcm.minApprovals {
		commitMsg, err := NewMessage(cmsg, types.ConsensusMessageTypeCommit, pcm.privKey)
		if err != nil {
			logrus.Errorf("failed to create commit message: %v", err)
		}
		pcm.psb.BroadcastToServiceTopic(commitMsg)
		pcm.state.status = StateStatusPrepared
	}
}

func (pcm *PBFTConsensusManager) handleCommit(message *pubsub.GenericMessage) {
	pcm.state.mutex.Lock()
	defer pcm.state.mutex.Unlock()
	commit, ok := message.Payload.(types.CommitMessage)
	if !ok {
		logrus.Warn("failed to convert payload to Prepare message")
		return
	}

	cmsg := types.ConsensusMessage{
		Type:      types.ConsensusMessageTypeCommit,
		From:      message.From,
		Blockhash: commit.Blockhash,
		Signature: commit.Signature, // TODO check the signature
	}

	if _, err := pcm.blockPool.GetBlock(cmsg.Blockhash); errors.Is(err, cache.ErrNotFound) {
		logrus.Debugf("received unknown block")
		return
	}

	if pcm.msgLog.Exists(cmsg) {
		logrus.Debugf("received existing commit msg, dropping...")
		return
	}
	if !pcm.validator.Valid(cmsg, nil) {
		logrus.Warn("received invalid commit msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(cmsg)

	if len(pcm.msgLog.Get(types.ConsensusMessageTypeCommit, cmsg.Blockhash)) >= pcm.minApprovals {
		block, err := pcm.blockPool.GetBlock(cmsg.Blockhash)
		if err != nil {
			logrus.Debug(err)
			return
		}
		pcm.blockPool.AddAcceptedBlock(block)
		pcm.state.status = StateStatusCommited
	}
}

func (pcm *PBFTConsensusManager) NewDrandRound(from phony.Actor, res client.Result) {
	pcm.Act(from, func() {
		pcm.state.mutex.Lock()
		defer pcm.state.mutex.Unlock()
		block, err := pcm.commitAcceptedBlocks()
		if err != nil {
			if errors.Is(err, ErrNoAcceptedBlocks) {
				logrus.Warnf("No accepted blocks for consensus round %d", pcm.state.blockHeight)
			} else {
				logrus.Errorf("Failed to select the block in consensus round %d: %s", pcm.state.blockHeight, err.Error())
			}
			return
		}

		// if we are miner for this block
		// then post dione tasks to target chains (currently, only Ethereum)
		if block.Header.Proposer == pcm.miner.address {
			for _, v := range block.Data {
				var task types2.DioneTask
				err := cbor.Unmarshal(v.Data, &task)
				if err != nil {
					logrus.Errorf("Failed to unmarshal transaction %x payload: %s", v.Hash, err.Error())
					continue // FIXME
				}
				reqIDNumber, ok := big.NewInt(0).SetString(task.RequestID, 10)
				if !ok {
					logrus.Errorf("Failed to parse request id number in task of tx %x", v.Hash)
					continue // FIXME
				}

				err = pcm.ethereumClient.SubmitRequestAnswer(reqIDNumber, task.Payload)
				if err != nil {
					logrus.Errorf("Failed to submit task in tx %x: %s", v.Hash, err.Error())
					continue // FIXME
				}
			}
		}

		pcm.state.ready <- true

		minedBlock, err := pcm.miner.MineBlock(res.Randomness(), block.Header)
		if err != nil {
			logrus.Errorf("Failed to mine the block: %s", err.Error())
			return
		}

		pcm.state.drandRound = res.Round()
		pcm.state.randomness = res.Randomness()
		pcm.state.blockHeight = pcm.state.blockHeight + 1

		// if we are round winner
		if minedBlock != nil {
			err = pcm.propose(minedBlock)
			if err != nil {
				logrus.Errorf("Failed to propose the block: %s", err.Error())
				return
			}
		}
	})
}

func (pcm *PBFTConsensusManager) commitAcceptedBlocks() (*types3.Block, error) {
	blocks := pcm.blockPool.GetAllAcceptedBlocks()
	if blocks == nil {
		return nil, ErrNoAcceptedBlocks
	}
	var maxStake *big.Int
	var selectedBlock *types3.Block
	for _, v := range blocks {
		stake, err := pcm.ethereumClient.GetMinerStake(v.Header.ProposerEth)
		if err != nil {
			return nil, err
		}

		if maxStake != nil {
			if stake.Cmp(maxStake) == -1 {
				continue
			}
		}
		maxStake = stake
		selectedBlock = v
	}
	logrus.Debugf("Selected block of miner %s", selectedBlock.Header.ProposerEth.Hex())
	pcm.blockPool.PruneAcceptedBlocks()
	return selectedBlock, pcm.blockchain.StoreBlock(selectedBlock)
}
