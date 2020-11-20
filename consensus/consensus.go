package consensus

import (
	"math/big"
	"sync"

	"github.com/Secured-Finance/dione/sigs"

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
}

type ConsensusData struct {
	//	preparedCount             int
	//	commitCount               int
	mutex            sync.Mutex
	alreadySubmitted bool
	//	result                    string
	//	test                      bool
	//	onConsensusFinishCallback func(finalData string)
}

func NewPBFTConsensusManager(psb *pubsub.PubSubRouter, minApprovals int, privKey []byte, ethereumClient *ethclient.EthereumClient) *PBFTConsensusManager {
	pcm := &PBFTConsensusManager{}
	pcm.psb = psb
	pcm.prePreparePool = NewPrePreparePool()
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

//func (pcm *PBFTConsensusManager) NewTestConsensus(data string, consensusID string, onConsensusFinishCallback func(finalData string)) {
//	//consensusID := uuid.New().String()
//	cData := &ConsensusData{}
//	cData.test = true
//	cData.onConsensusFinishCallback = onConsensusFinishCallback
//	pcm.Consensuses[consensusID] = cData
//
//	// here we will create DioneTask
//
//	msg := models.Message{}
//	msg.Type = "prepared"
//	msg.Payload = make(map[string]interface{})
//	msg.Payload["consensusID"] = consensusID
//	msg.Payload["data"] = data
//	sign, err := sigs.Sign(types.SigTypeEd25519, pcm.privKey, []byte(data))
//	if err != nil {
//		logrus.Warnf("failed to sign data: %w", err)
//		return
//	}
//	msg.Payload["signature"] = string(sign.Data)
//	pcm.psb.BroadcastToServiceTopic(&msg)
//
//	cData.State = ConsensusPrePrepared
//	logrus.Debug("started new consensus: " + consensusID)
//}
//
//func (pcm *PBFTConsensusManager) handlePrePrepareMessage(sender peer.ID, message *models.Message) {
//	consensusID := message.Payload["consensusID"].(string)
//	if _, ok := pcm.Consensuses[consensusID]; !ok {
//		logrus.Warn("Unknown consensus ID,: " + consensusID + ", creating consensusInfo")
//		pcm.Consensuses[consensusID] = &ConsensusData{
//			State:                     ConsensusPrePrepared,
//			onConsensusFinishCallback: func(finalData string) {},
//		}
//	}
//	data := pcm.Consensuses[consensusID]
//	logrus.Debug("received pre_prepare msg")
//}
//
//func (pcm *PBFTConsensusManager) handlePreparedMessage(sender peer.ID, message *models.Message) {
//	// TODO add check on view of the message
//	consensusID := message.Payload["consensusID"].(string)
//	if _, ok := pcm.Consensuses[consensusID]; !ok {
//		logrus.Warn("Unknown consensus ID,: " + consensusID + ", creating consensusInfo")
//		pcm.Consensuses[consensusID] = &ConsensusData{
//			State:                     ConsensusPrePrepared,
//			onConsensusFinishCallback: func(finalData string) {},
//		}
//	}
//	data := pcm.Consensuses[consensusID]
//	logrus.Debug("received prepared msg")
//	//data := pcm.Consensuses[consensusID]
//
//	// TODO
//	// here we can validate miner which produced this task, is he winner, and so on
//	// validation steps:
//	// 1. validate sender eligibility to mine (check if it has minimal stake)
//	// 2. validate sender wincount
//	// 3. validate randomness
//	// 4. validate vrf
//	// 5. validate payload signature
//	// 6. validate transaction (get from rpc client and compare with received)
//
//	signStr := message.Payload["signature"].(string)
//	signRaw := []byte(signStr)
//	err := sigs.Verify(&types.Signature{Data: signRaw, Type: types.SigTypeEd25519}, sender, message.Payload["data"].([]byte))
//	if err != nil {
//		logrus.Warn("failed to verify data signature")
//		return
//	}
//
//	data.mutex.Lock()
//	defer data.mutex.Unlock()
//	data.preparedCount++
//
//	if data.preparedCount > 2*pcm.maxFaultNodes+1 {
//		msg := models.Message{}
//		msg.Payload = make(map[string]interface{})
//		msg.Type = "commit"
//		msg.Payload["consensusID"] = consensusID
//		msg.Payload["data"] = message.Payload["data"]
//		sign, err := sigs.Sign(types.SigTypeEd25519, pcm.privKey, message.Payload["data"].([]byte))
//		if err != nil {
//			logrus.Warnf("failed to sign data: %w", err)
//			return
//		}
//		msg.Payload["signature"] = string(sign.Data)
//		err = pcm.psb.BroadcastToServiceTopic(&msg)
//		if err != nil {
//			logrus.Warn("Unable to send COMMIT message: " + err.Error())
//			return
//		}
//		data.State = ConsensusPrepared
//	}
//}
//
//func (pcm *PBFTConsensusManager) handleCommitMessage(sender peer.ID, message *models.Message) {
//	// TODO add check on view of the message
//	// TODO add validation of data by hash to this stage
//	consensusID := message.Payload["consensusID"].(string)
//	if _, ok := pcm.Consensuses[consensusID]; !ok {
//		logrus.Warn("Unknown consensus ID: " + consensusID)
//		return
//	}
//	data := pcm.Consensuses[consensusID]
//
//	data.mutex.Lock()
//	defer data.mutex.Unlock()
//	if data.State == ConsensusCommitted {
//		logrus.Debug("consensus already finished, dropping COMMIT message")
//		return
//	}
//
//	logrus.Debug("received commit msg")
//
//	signStr := message.Payload["signature"].(string)
//	signRaw := []byte(signStr)
//	err := sigs.Verify(&types.Signature{Data: signRaw, Type: types.SigTypeEd25519}, sender, message.Payload["data"].([]byte))
//	if err != nil {
//		logrus.Warn("failed to verify data signature")
//		return
//	}
//
//	data.commitCount++
//
//	if data.commitCount > 2*pcm.maxFaultNodes+1 {
//		data.State = ConsensusCommitted
//		data.result = message.Payload["data"].(string)
//		logrus.Debug("consensus successfully finished with result: " + data.result)
//		data.onConsensusFinishCallback(data.result)
//	}
//}

func (pcm *PBFTConsensusManager) Propose(consensusID, data string, requestID *big.Int, callbackAddress common.Address) error {
	pcm.consensusInfo[consensusID] = &ConsensusData{}
	reqIDRaw := requestID.String()
	callbackAddressHex := callbackAddress.Hex()
	prePrepareMsg, err := pcm.prePreparePool.CreatePrePrepare(consensusID, data, reqIDRaw, callbackAddressHex, pcm.privKey)
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

	err := pcm.resignMessage(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	err = pcm.psb.BroadcastToServiceTopic(message)
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
	err := pcm.resignMessage(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	err = pcm.psb.BroadcastToServiceTopic(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	if pcm.preparePool.PrepareSize(message.Payload.ConsensusID) >= pcm.minApprovals {
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
	err := pcm.resignMessage(message)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	err = pcm.psb.BroadcastToServiceTopic(message)
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
			callbackAddress := common.HexToAddress(consensusMsg.CallbackAddress)
			err := pcm.ethereumClient.SubmitRequestAnswer(reqID, consensusMsg.Data, callbackAddress)
			if err != nil {
				logrus.Errorf("Failed to submit on-chain result: %w", err)
			}
			info.alreadySubmitted = true
		}
	}
}

func (pcm *PBFTConsensusManager) resignMessage(msg *types.Message) error {
	sig, err := sigs.Sign(types2.SigTypeEd25519, pcm.privKey, []byte(msg.Payload.Data))
	if err != nil {
		return err
	}
	msg.Payload.Signature = sig.Data
	return nil
}
