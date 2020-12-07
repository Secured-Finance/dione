package consensus

import (
	"sync"

	types2 "github.com/Secured-Finance/dione/consensus/types"
)

type PreparePool struct {
	mut         sync.RWMutex
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
	newCMsg := prePrepareMsg.Payload
	message.Payload = newCMsg
	return &message, nil
}

func (pp *PreparePool) IsExistingPrepare(prepareMsg *types2.Message) bool {
	pp.mut.RLock()
	defer pp.mut.RUnlock()

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
	err := verifyTaskSignature(consensusMsg)
	if err != nil {
		return false
	}
	return true
}

func (pp *PreparePool) AddPrepare(prepare *types2.Message) {
	pp.mut.Lock()
	defer pp.mut.Unlock()

	consensusID := prepare.Payload.ConsensusID
	if _, ok := pp.prepareMsgs[consensusID]; !ok {
		pp.prepareMsgs[consensusID] = []*types2.Message{}
	}

	pp.prepareMsgs[consensusID] = append(pp.prepareMsgs[consensusID], prepare)
}

func (pp *PreparePool) PreparePoolSize(consensusID string) int {
	pp.mut.RLock()
	defer pp.mut.RUnlock()

	if v, ok := pp.prepareMsgs[consensusID]; ok {
		return len(v)
	}
	return 0
}
