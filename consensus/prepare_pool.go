package consensus

import (
	"github.com/Secured-Finance/dione/models"
	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
)

type PreparePool struct {
	prepareMsgs map[string][]*models.Message
	privateKey  []byte
}

func NewPreparePool() *PreparePool {
	return &PreparePool{
		prepareMsgs: map[string][]*models.Message{},
	}
}

func (pp *PreparePool) CreatePrepare(prePrepareMsg *models.Message, privateKey []byte) (*models.Message, error) {
	var message models.Message
	message.Type = models.MessageTypePrepare
	var consensusMsg models.ConsensusMessage
	prepareCMessage := prePrepareMsg.Payload
	consensusMsg.ConsensusID = prepareCMessage.ConsensusID
	consensusMsg.RequestID = prePrepareMsg.Payload.RequestID
	consensusMsg.CallbackAddress = prePrepareMsg.Payload.CallbackAddress
	consensusMsg.Data = prepareCMessage.Data
	signature, err := sigs.Sign(types.SigTypeEd25519, privateKey, []byte(prepareCMessage.Data))
	if err != nil {
		return nil, err
	}
	consensusMsg.Signature = signature.Data
	message.Payload = consensusMsg
	return &message, nil
}

func (pp *PreparePool) IsExistingPrepare(prepareMsg *models.Message) bool {
	consensusMessage := prepareMsg.Payload
	var exists bool
	for _, v := range pp.prepareMsgs[consensusMessage.ConsensusID] {
		if v.From == prepareMsg.From {
			exists = true
		}
	}
	return exists
}

func (pp *PreparePool) IsValidPrepare(prepare *models.Message) bool {
	consensusMsg := prepare.Payload
	err := sigs.Verify(&types.Signature{Type: types.SigTypeEd25519, Data: consensusMsg.Signature}, prepare.From, []byte(consensusMsg.Data))
	if err != nil {
		return false
	}
	return true
}

func (pp *PreparePool) AddPrepare(prepare *models.Message) {
	consensusID := prepare.Payload.ConsensusID
	if _, ok := pp.prepareMsgs[consensusID]; !ok {
		pp.prepareMsgs[consensusID] = []*models.Message{}
	}

	pp.prepareMsgs[consensusID] = append(pp.prepareMsgs[consensusID], prepare)
}

func (pp *PreparePool) PrepareSize(consensusID string) int {
	if v, ok := pp.prepareMsgs[consensusID]; ok {
		return len(v)
	}
	return 0
}
