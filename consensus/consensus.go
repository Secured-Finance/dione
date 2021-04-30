package consensus

import (
	"math/big"
	"sync"

	"github.com/Secured-Finance/dione/cache"

	"github.com/Secured-Finance/dione/consensus/types"

	"github.com/Secured-Finance/dione/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/Secured-Finance/dione/pubsub"
	types2 "github.com/Secured-Finance/dione/types"
)

type PBFTConsensusManager struct {
	psb            *pubsub.PubSubRouter
	minApprovals   int
	privKey        []byte
	msgLog         *MessageLog
	validator      *ConsensusValidator
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
	pcm.validator = NewConsensusValidator(evc, miner)
	pcm.msgLog = NewMessageLog()
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

func (pcm *PBFTConsensusManager) Propose(task types2.DioneTask) error {
	pcm.createConsensusInfo(&task, true)

	prePrepareMsg, err := CreatePrePrepareWithTaskSignature(&task, pcm.privKey)
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
	if pcm.msgLog.Exists(*message) {
		return
	}
	if !pcm.validator.Valid(*message) {
		logrus.Warn("received invalid pre_prepare msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(*message)
	err := pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	prepareMsg, err := NewMessage(message, types.MessageTypePrepare)
	if err != nil {
		logrus.Errorf("failed to create prepare message: %w", err)
	}

	pcm.createConsensusInfo(&message.Payload.Task, false)

	pcm.psb.BroadcastToServiceTopic(&prepareMsg)
}

func (pcm *PBFTConsensusManager) handlePrepare(message *types.Message) {
	if pcm.msgLog.Exists(*message) {
		return
	}
	if !pcm.validator.Valid(*message) {
		logrus.Warn("received invalid prepare msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(*message)
	err := pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	if len(pcm.msgLog.GetMessagesByTypeAndConsensusID(types.MessageTypePrepare, message.Payload.Task.ConsensusID)) >= pcm.minApprovals {
		commitMsg, err := NewMessage(message, types.MessageTypeCommit)
		if err != nil {
			logrus.Errorf("failed to create commit message: %w", err)
		}
		pcm.psb.BroadcastToServiceTopic(&commitMsg)
	}
}

func (pcm *PBFTConsensusManager) handleCommit(message *types.Message) {
	if pcm.msgLog.Exists(*message) {
		return
	}
	if !pcm.validator.Valid(*message) {
		logrus.Warn("received invalid commit msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(*message)
	err := pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	consensusMsg := message.Payload
	if len(pcm.msgLog.GetMessagesByTypeAndConsensusID(types.MessageTypeCommit, message.Payload.Task.ConsensusID)) >= pcm.minApprovals {
		info := pcm.GetConsensusInfo(consensusMsg.Task.ConsensusID)
		if info == nil {
			logrus.Debugf("consensus doesn't exist in our consensus map - skipping...")
			return
		}
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

			err = pcm.ethereumClient.SubmitRequestAnswer(reqID, consensusMsg.Task.Payload)
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
