package consensus

import (
	"math/big"
	"sync"

	"github.com/Secured-Finance/dione/cache"

	oracleEmitter "github.com/Secured-Finance/dione/contracts/oracleemitter"

	"github.com/Secured-Finance/dione/consensus/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Secured-Finance/dione/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/Secured-Finance/dione/pubsub"
	types2 "github.com/Secured-Finance/dione/types"
)

type PBFTConsensusManager struct {
	psb            *pubsub.PubSubRouter
	minApprovals   int
	privKey        []byte
	prePreparePool *PrePreparePool
	preparePool    *PreparePool
	commitPool     *CommitPool
	consensusInfo  map[string]*ConsensusData
	ethereumClient *ethclient.EthereumClient
	miner          *Miner
}

type ConsensusData struct {
	mutex            sync.Mutex
	alreadySubmitted bool
}

func NewPBFTConsensusManager(psb *pubsub.PubSubRouter, minApprovals int, privKey []byte, ethereumClient *ethclient.EthereumClient, miner *Miner, evc *cache.EventLogCache) *PBFTConsensusManager {
	pcm := &PBFTConsensusManager{}
	pcm.psb = psb
	pcm.miner = miner
	pcm.prePreparePool = NewPrePreparePool(miner, evc)
	pcm.preparePool = NewPreparePool()
	pcm.commitPool = NewCommitPool()
	pcm.minApprovals = minApprovals
	pcm.privKey = privKey
	pcm.ethereumClient = ethereumClient
	pcm.consensusInfo = map[string]*ConsensusData{}
	pcm.psb.Hook(types.MessageTypePrePrepare, pcm.handlePrePrepare)
	pcm.psb.Hook(types.MessageTypePrepare, pcm.handlePrepare)
	pcm.psb.Hook(types.MessageTypeCommit, pcm.handleCommit)
	return pcm
}

func (pcm *PBFTConsensusManager) Propose(consensusID string, task types2.DioneTask, requestEvent *oracleEmitter.OracleEmitterNewOracleRequest) error {
	pcm.consensusInfo[consensusID] = &ConsensusData{}

	prePrepareMsg, err := pcm.prePreparePool.CreatePrePrepare(
		consensusID,
		task,
		requestEvent.RequestID.String(),
		requestEvent.CallbackAddress.Bytes(),
		requestEvent.CallbackMethodID[:],
		pcm.privKey,
	)
	if err != nil {
		return err
	}
	pcm.psb.BroadcastToServiceTopic(prePrepareMsg)
	return nil
}

func (pcm *PBFTConsensusManager) handlePrePrepare(message *types.Message) {
	if pcm.prePreparePool.IsExistingPrePrepare(message) {
		logrus.Debug("received existing pre_prepare msg, dropping...")
		return
	}
	if !pcm.prePreparePool.IsValidPrePrepare(message) {
		logrus.Debug("received invalid pre_prepare msg, dropping...")
		return
	}

	pcm.prePreparePool.AddPrePrepare(message)
	err := pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	prepareMsg, err := pcm.preparePool.CreatePrepare(message, pcm.privKey)
	if err != nil {
		logrus.Errorf("failed to create prepare message: %w", err)
	}
	pcm.psb.BroadcastToServiceTopic(prepareMsg)
}

func (pcm *PBFTConsensusManager) handlePrepare(message *types.Message) {
	if pcm.preparePool.IsExistingPrepare(message) {
		logrus.Debug("received existing prepare msg, dropping...")
		return
	}
	if !pcm.preparePool.IsValidPrepare(message) {
		logrus.Debug("received invalid prepare msg, dropping...")
		return
	}

	pcm.preparePool.AddPrepare(message)
	err := pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	if pcm.preparePool.PreparePoolSize(message.Payload.ConsensusID) >= pcm.minApprovals {
		commitMsg, err := pcm.commitPool.CreateCommit(message, pcm.privKey)
		if err != nil {
			logrus.Errorf("failed to create commit message: %w", err)
		}
		pcm.psb.BroadcastToServiceTopic(commitMsg)
	}
}

func (pcm *PBFTConsensusManager) handleCommit(message *types.Message) {
	if pcm.commitPool.IsExistingCommit(message) {
		logrus.Debug("received existing commit msg, dropping...")
		return
	}
	if !pcm.commitPool.IsValidCommit(message) {
		logrus.Debug("received invalid commit msg, dropping...")
		return
	}

	pcm.commitPool.AddCommit(message)
	err := pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	consensusMsg := message.Payload
	if pcm.commitPool.CommitSize(consensusMsg.ConsensusID) >= pcm.minApprovals {
		if info, ok := pcm.consensusInfo[consensusMsg.ConsensusID]; ok {
			info.mutex.Lock()
			defer info.mutex.Unlock()
			if info.alreadySubmitted {
				return
			}
			logrus.Infof("Submitting on-chain result for consensus ID: %s", consensusMsg.ConsensusID)
			reqID, ok := new(big.Int).SetString(consensusMsg.RequestID, 10)
			if !ok {
				logrus.Errorf("Failed to parse big int: %v", consensusMsg.RequestID)
			}
			callbackAddress := common.BytesToAddress(consensusMsg.CallbackAddress)
			err := pcm.ethereumClient.SubmitRequestAnswer(reqID, string(consensusMsg.Task.Payload), callbackAddress)
			if err != nil {
				logrus.Errorf("Failed to submit on-chain result: %w", err)
			}
			info.alreadySubmitted = true
		}
	}
}
