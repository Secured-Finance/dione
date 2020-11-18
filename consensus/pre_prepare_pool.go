package consensus

import (
	"github.com/Secured-Finance/dione/models"
	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
	"github.com/sirupsen/logrus"
)

type PrePreparePool struct {
	prePrepareMsgs map[string][]*models.Message
}

func NewPrePreparePool() *PrePreparePool {
	return &PrePreparePool{
		prePrepareMsgs: map[string][]*models.Message{},
	}
}

func (pp *PrePreparePool) CreatePrePrepare(consensusID, data string, requestID string, callbackAddress string, privateKey []byte) (*models.Message, error) {
	var message models.Message
	message.Type = models.MessageTypePrePrepare
	var consensusMsg models.ConsensusMessage
	consensusMsg.ConsensusID = consensusID
	consensusMsg.RequestID = requestID
	consensusMsg.CallbackAddress = callbackAddress
	consensusMsg.Data = data
	signature, err := sigs.Sign(types.SigTypeEd25519, privateKey, []byte(data))
	if err != nil {
		return nil, err
	}
	consensusMsg.Signature = signature.Data
	message.Payload = consensusMsg
	return &message, nil
}

func (ppp *PrePreparePool) IsExistingPrePrepare(prepareMsg *models.Message) bool {
	consensusMessage := prepareMsg.Payload
	var exists bool
	for _, v := range ppp.prePrepareMsgs[consensusMessage.ConsensusID] {
		if v.From == prepareMsg.From {
			exists = true
		}
	}
	return exists
}

func (ppp *PrePreparePool) IsValidPrePrepare(prePrepare *models.Message) bool {
	// TODO here we need to do validation of tx itself
	consensusMsg := prePrepare.Payload
	err := sigs.Verify(&types.Signature{Type: types.SigTypeEd25519, Data: consensusMsg.Signature}, prePrepare.From, []byte(consensusMsg.Data))
	if err != nil {
		logrus.Errorf("unable to verify signature: %v", err)
		return false
	}
	return true
}

func (ppp *PrePreparePool) AddPrePrepare(prePrepare *models.Message) {
	consensusID := prePrepare.Payload.ConsensusID
	if _, ok := ppp.prePrepareMsgs[consensusID]; !ok {
		ppp.prePrepareMsgs[consensusID] = []*models.Message{}
	}

	ppp.prePrepareMsgs[consensusID] = append(ppp.prePrepareMsgs[consensusID], prePrepare)
}
