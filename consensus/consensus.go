package consensus

import (
	"math/big"
	"sync"

	"github.com/fxamacker/cbor/v2"

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
	cache          cache.Cache
}

type Consensus struct {
	mutex                sync.Mutex
	Finished             bool
	IsCurrentMinerLeader bool
	Task                 *types2.DioneTask
}

func NewPBFTConsensusManager(psb *pubsub.PubSubRouter, minApprovals int, privKey []byte, ethereumClient *ethclient.EthereumClient, miner *Miner, evc cache.Cache) *PBFTConsensusManager {
	pcm := &PBFTConsensusManager{}
	pcm.psb = psb
	pcm.miner = miner
	pcm.validator = NewConsensusValidator(evc, miner)
	pcm.msgLog = NewMessageLog()
	pcm.minApprovals = minApprovals
	pcm.privKey = privKey
	pcm.ethereumClient = ethereumClient
	pcm.cache = evc
	pcm.consensusMap = map[string]*Consensus{}
	pcm.psb.Hook(pubsub.PrePrepareMessageType, pcm.handlePrePrepare)
	pcm.psb.Hook(pubsub.PrepareMessageType, pcm.handlePrepare)
	pcm.psb.Hook(pubsub.CommitMessageType, pcm.handleCommit)
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

func (pcm *PBFTConsensusManager) handlePrePrepare(message *pubsub.PubSubMessage) {
	cmsg, err := unmarshalPayload(message)
	if err != nil {
		return
	}

	if cmsg.Task.Miner == pcm.miner.address {
		return
	}
	if pcm.msgLog.Exists(cmsg) {
		logrus.Debugf("received existing pre_prepare msg, dropping...")
		return
	}
	if !pcm.validator.Valid(cmsg) {
		logrus.Warn("received invalid pre_prepare msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(cmsg)

	prepareMsg, err := NewMessage(message, pubsub.PrepareMessageType)
	if err != nil {
		logrus.Errorf("failed to create prepare message: %v", err)
	}

	pcm.createConsensusInfo(&cmsg.Task, false)

	pcm.psb.BroadcastToServiceTopic(&prepareMsg)
}

func (pcm *PBFTConsensusManager) handlePrepare(message *pubsub.PubSubMessage) {
	cmsg, err := unmarshalPayload(message)
	if err != nil {
		return
	}

	if pcm.msgLog.Exists(cmsg) {
		logrus.Debugf("received existing prepare msg, dropping...")
		return
	}
	if !pcm.validator.Valid(cmsg) {
		logrus.Warn("received invalid prepare msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(cmsg)

	if len(pcm.msgLog.Get(types.MessageTypePrepare, cmsg.Task.ConsensusID)) >= pcm.minApprovals {
		commitMsg, err := NewMessage(message, types.MessageTypeCommit)
		if err != nil {
			logrus.Errorf("failed to create commit message: %w", err)
		}
		pcm.psb.BroadcastToServiceTopic(&commitMsg)
	}
}

func (pcm *PBFTConsensusManager) handleCommit(message *pubsub.PubSubMessage) {
	cmsg, err := unmarshalPayload(message)
	if err != nil {
		return
	}

	if pcm.msgLog.Exists(cmsg) {
		logrus.Debugf("received existing commit msg, dropping...")
		return
	}
	if !pcm.validator.Valid(cmsg) {
		logrus.Warn("received invalid commit msg, dropping...")
		return
	}

	pcm.msgLog.AddMessage(cmsg)

	if len(pcm.msgLog.Get(types.MessageTypeCommit, cmsg.Task.ConsensusID)) >= pcm.minApprovals {
		info := pcm.GetConsensusInfo(cmsg.Task.ConsensusID)
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
			logrus.Infof("Submitting on-chain result for consensus ID: %s", cmsg.Task.ConsensusID)
			reqID, ok := new(big.Int).SetString(cmsg.Task.RequestID, 10)
			if !ok {
				logrus.Errorf("Failed to parse request ID: %v", cmsg.Task.RequestID)
			}

			err := pcm.ethereumClient.SubmitRequestAnswer(reqID, cmsg.Task.Payload)
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

func unmarshalPayload(msg *pubsub.PubSubMessage) (types.ConsensusMessage, error) {
	var task types2.DioneTask
	err := cbor.Unmarshal(msg.Payload, &task)
	if err != nil {
		logrus.Debug(err)
		return types.ConsensusMessage{}, err
	}
	var consensusMessageType types.MessageType
	switch msg.Type {
	case pubsub.PrePrepareMessageType:
		{
			consensusMessageType = types.MessageTypePrePrepare
			break
		}
	case pubsub.PrepareMessageType:
		{
			consensusMessageType = types.MessageTypePrepare
			break
		}
	case pubsub.CommitMessageType:
		{
			consensusMessageType = types.MessageTypeCommit
			break
		}
	}
	cmsg := types.ConsensusMessage{
		Type: consensusMessageType,
		From: msg.From,
		Task: task,
	}
	return cmsg, nil
}
