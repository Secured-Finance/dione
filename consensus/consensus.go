package consensus

import (
	"sync"

	"github.com/Secured-Finance/dione/models"
	"github.com/Secured-Finance/dione/pb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ConsensusState int

const (
	consensusPrePrepared ConsensusState = 0x0
	consensusPrepared    ConsensusState = 0x1
	consensusCommitted   ConsensusState = 0x2

	testValidData = "test"
)

type PBFTConsensusManager struct {
	psb           *pb.PubSubRouter
	Consensuses   map[string]*ConsensusData
	maxFaultNodes int
}

type ConsensusData struct {
	preparedCount int
	commitCount   int
	State         ConsensusState
	mutex         sync.Mutex
	test          bool
}

func NewPBFTConsensusManager(psb *pb.PubSubRouter, maxFaultNodes int) *PBFTConsensusManager {
	pcm := &PBFTConsensusManager{}
	pcm.Consensuses = make(map[string]*ConsensusData)
	pcm.psb = psb
	pcm.psb.Hook("prepared", pcm.handlePreparedMessage)
	pcm.psb.Hook("commit", pcm.handleCommitMessage)
	return pcm
}

func (pcm *PBFTConsensusManager) NewTestConsensus(data string) {
	consensusID := uuid.New().String()
	cData := &ConsensusData{}
	cData.test = true
	pcm.Consensuses[consensusID] = cData

	msg := models.Message{}
	msg.Type = "prepared"
	msg.Payload = make(map[string]interface{})
	msg.Payload["consensusID"] = consensusID
	msg.Payload["data"] = data
	pcm.psb.BroadcastToServiceTopic(&msg)

	cData.State = consensusPrePrepared
	logrus.Debug("started new consensus: " + consensusID)
}

func (pcm *PBFTConsensusManager) handlePreparedMessage(message *models.Message) {
	// TODO add check on view of the message
	consensusID := message.Payload["consensusID"].(string)
	if _, ok := pcm.Consensuses[consensusID]; !ok {
		logrus.Warn("Unknown consensus ID: " + consensusID)
		return
	}
	data := pcm.Consensuses[consensusID]

	// validate payload data
	if data.test {
		rData := message.Payload["data"].(string)
		if rData != testValidData {
			logrus.Error("Incorrect data was received! Ignoring this message, because it was sent from fault node!")
			return
		}
	} else {
		// TODO
	}

	data.mutex.Lock()
	data.preparedCount++
	data.mutex.Unlock()

	if data.preparedCount > 2*pcm.maxFaultNodes+1 {
		msg := models.Message{}
		msg.Type = "commit"
		msg.Payload["consensusID"] = consensusID
		err := pcm.psb.BroadcastToServiceTopic(&msg)
		if err != nil {
			logrus.Warn("Unable to send COMMIT message: " + err.Error())
			return
		}
		data.State = consensusPrepared
	}
}

func (pcm *PBFTConsensusManager) handleCommitMessage(message *models.Message) {
	// TODO add check on view of the message
	consensusID := message.Payload["consensusID"].(string)
	if _, ok := pcm.Consensuses[consensusID]; !ok {
		logrus.Warn("Unknown consensus ID: " + consensusID)
		return
	}
	data := pcm.Consensuses[consensusID]
	data.mutex.Lock()
	data.commitCount++
	data.mutex.Unlock()

	if data.commitCount > 2*pcm.maxFaultNodes+1 {
		logrus.Debug("consensus successfully finished")
		data.State = consensusCommitted
	}
}
