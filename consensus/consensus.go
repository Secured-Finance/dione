package consensus

import (
	"math/big"
	"sync"

	"github.com/Secured-Finance/dione/cache"

	"github.com/Secured-Finance/dione/contracts/dioneOracle"

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
	consensusMap   map[string]*Consensus
	ethereumClient *ethclient.EthereumClient
	miner          *Miner
	eventCache     cache.EventCache
}

type Consensus struct {
	mutex                sync.Mutex
	Finished             bool
	IsCurrentMinerLeader bool
	Task                 *types2.DioneTask
}

func NewPBFTConsensusManager(psb *pubsub.PubSubRouter, minApprovals int, privKey []byte, ethereumClient *ethclient.EthereumClient, miner *Miner, evc cache.EventCache) *PBFTConsensusManager {
	pcm := &PBFTConsensusManager{}
	pcm.psb = psb
	pcm.miner = miner
	pcm.prePreparePool = NewPrePreparePool(miner, evc)
	pcm.preparePool = NewPreparePool()
	pcm.commitPool = NewCommitPool()
	pcm.minApprovals = minApprovals
	pcm.privKey = privKey
	pcm.ethereumClient = ethereumClient
	pcm.eventCache = evc
	pcm.consensusMap = map[string]*Consensus{}
	pcm.psb.Hook(types.MessageTypePrePrepare, pcm.handlePrePrepare)
	pcm.psb.Hook(types.MessageTypePrepare, pcm.handlePrepare)
	pcm.psb.Hook(types.MessageTypeCommit, pcm.handleCommit)
	return pcm
}

func (pcm *PBFTConsensusManager) Propose(consensusID string, task types2.DioneTask, requestEvent *dioneOracle.DioneOracleNewOracleRequest) error {
	pcm.createConsensusInfo(&task, true)

	prePrepareMsg, err := pcm.prePreparePool.CreatePrePrepare(
		&task,
		requestEvent.ReqID.String(),
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
	if message.Payload.Task.Miner == pcm.miner.address {
		return
	}
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

	pcm.createConsensusInfo(&message.Payload.Task, false)

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

	if pcm.preparePool.PreparePoolSize(message.Payload.Task.ConsensusID) >= pcm.minApprovals {
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
	if pcm.commitPool.CommitSize(consensusMsg.Task.ConsensusID) >= pcm.minApprovals {
		info := pcm.consensusMap[consensusMsg.Task.ConsensusID]
		info.mutex.Lock()
		defer info.mutex.Unlock()
		if info.Finished {
			return
		}
		if info.IsCurrentMinerLeader {
			logrus.Infof("Submitting on-chain result for consensus ID: %s", consensusMsg.Task.ConsensusID)
			reqID, ok := new(big.Int).SetString(consensusMsg.Task.RequestID, 10)
			if !ok {
				logrus.Errorf("Failed to parse request ID: %v", consensusMsg.Task.RequestID)
			}
			callbackAddress := common.BytesToAddress(consensusMsg.Task.CallbackAddress)

			request, err := pcm.eventCache.GetOracleRequestEvent("request_" + consensusMsg.Task.RequestID)
			if err != nil {
				logrus.Errorf("Failed to get request from cache: %v", err.Error())
				return
			}

			err = pcm.ethereumClient.SubmitRequestAnswer(reqID, callbackAddress, request.CallbackMethodID, request.RequestParams, request.Deadline, consensusMsg.Task.Payload)
			if err != nil {
				logrus.Errorf("Failed to submit on-chain result: %v", err)
			}
		}

		info.Finished = true
	}
}

func (pcm *PBFTConsensusManager) createConsensusInfo(task *types2.DioneTask, isLeader bool) {
	if _, ok := pcm.consensusMap[task.ConsensusID]; !ok {
		pcm.consensusMap[task.ConsensusID] = &Consensus{
			IsCurrentMinerLeader: isLeader,
			Task:                 task,
			Finished:             false,
		}
	}
}

func (pcm *PBFTConsensusManager) GetConsensusInfo(consensusID string) *Consensus {
	c, ok := pcm.consensusMap[consensusID]
	if !ok {
		return nil
	}

	return c
}
