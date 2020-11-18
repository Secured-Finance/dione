package consensus

import (
	types2 "github.com/Secured-Finance/dione/consensus/types"
	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
)

type PreparePool struct {
	prepareMsgs map[string][]*types2.Message
	privateKey  []byte
}

func NewPreparePool() *PreparePool {
	return &PreparePool{
		prepareMsgs: map[string][]*types2.Message{},
	}
}

func (pp *PreparePool) CreatePrepare(prePrepareMsg *types2.Message, privateKey []byte) (*types2.Message, error) {
	var message types2.Message
	message.Type = types2.MessageTypePrepare
	var consensusMsg types2.ConsensusMessage
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

func (pp *PreparePool) IsExistingPrepare(prepareMsg *types2.Message) bool {
	consensusMessage := prepareMsg.Payload
	var exists bool
	for _, v := range pp.prepareMsgs[consensusMessage.ConsensusID] {
		if v.From == prepareMsg.From {
			exists = true
		}
	}
	return exists
}

func (pp *PreparePool) IsValidPrepare(prepare *types2.Message) bool {
	consensusMsg := prepare.Payload
	err := sigs.Verify(&types.Signature{Type: types.SigTypeEd25519, Data: consensusMsg.Signature}, prepare.From, []byte(consensusMsg.Data))
	if err != nil {
		return false
	}
	return true
}

func (pp *PreparePool) AddPrepare(prepare *types2.Message) {
	consensusID := prepare.Payload.ConsensusID
	if _, ok := pp.prepareMsgs[consensusID]; !ok {
		pp.prepareMsgs[consensusID] = []*types2.Message{}
	}

	pp.prepareMsgs[consensusID] = append(pp.prepareMsgs[consensusID], prepare)
}

func (pp *PreparePool) PrepareSize(consensusID string) int {
	if v, ok := pp.prepareMsgs[consensusID]; ok {
		return len(v)
	}
	return 0
}
