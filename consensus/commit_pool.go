package consensus

import (
	types2 "github.com/Secured-Finance/dione/consensus/types"
	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
)

type CommitPool struct {
	commitMsgs map[string][]*types2.Message
}

func NewCommitPool() *CommitPool {
	return &CommitPool{
		commitMsgs: map[string][]*types2.Message{},
	}
}

func (cp *CommitPool) CreateCommit(prepareMsg *types2.Message, privateKey []byte) (*types2.Message, error) {
	var message types2.Message
	message.Type = types2.MessageTypeCommit
	var consensusMsg types2.ConsensusMessage
	prepareCMessage := prepareMsg.Payload
	consensusMsg.ConsensusID = prepareCMessage.ConsensusID
	consensusMsg.RequestID = prepareMsg.Payload.RequestID
	consensusMsg.CallbackAddress = prepareMsg.Payload.CallbackAddress
	consensusMsg.Data = prepareCMessage.Data
	signature, err := sigs.Sign(types.SigTypeEd25519, privateKey, []byte(prepareCMessage.Data))
	if err != nil {
		return nil, err
	}
	consensusMsg.Signature = signature.Data
	message.Payload = consensusMsg
	return &message, nil
}

func (cp *CommitPool) IsExistingCommit(commitMsg *types2.Message) bool {
	consensusMessage := commitMsg.Payload
	var exists bool
	for _, v := range cp.commitMsgs[consensusMessage.ConsensusID] {
		if v.From == commitMsg.From {
			exists = true
		}
	}
	return exists
}

func (cp *CommitPool) IsValidCommit(commit *types2.Message) bool {
	consensusMsg := commit.Payload
	err := sigs.Verify(&types.Signature{Type: types.SigTypeEd25519, Data: consensusMsg.Signature}, commit.From, []byte(consensusMsg.Data))
	if err != nil {
		return false
	}
	return true
}

func (cp *CommitPool) AddCommit(commit *types2.Message) {
	consensusID := commit.Payload.ConsensusID
	if _, ok := cp.commitMsgs[consensusID]; !ok {
		cp.commitMsgs[consensusID] = []*types2.Message{}
	}

	cp.commitMsgs[consensusID] = append(cp.commitMsgs[consensusID], commit)
}

func (cp *CommitPool) CommitSize(consensusID string) int {
	if v, ok := cp.commitMsgs[consensusID]; ok {
		return len(v)
	}
	return 0
}
