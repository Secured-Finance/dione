package consensus

import (
	"encoding/hex"

	"github.com/Secured-Finance/dione/models"
	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
)

type CommitPool struct {
	commitMsgs map[string][]*models.Message
}

func NewCommitPool() *CommitPool {
	return &CommitPool{
		commitMsgs: map[string][]*models.Message{},
	}
}

func (cp *CommitPool) CreateCommit(prepareMsg *models.Message, privateKey []byte) (*models.Message, error) {
	var message models.Message
	message.Type = models.MessageTypeCommit
	var consensusMsg models.ConsensusMessage
	prepareCMessage := prepareMsg.Payload
	consensusMsg.ConsensusID = prepareCMessage.ConsensusID
	consensusMsg.RequestID = prepareMsg.Payload.RequestID
	consensusMsg.CallbackAddress = prepareMsg.Payload.CallbackAddress
	consensusMsg.Data = prepareCMessage.Data
	signature, err := sigs.Sign(types.SigTypeEd25519, privateKey, []byte(prepareCMessage.Data))
	if err != nil {
		return nil, err
	}
	consensusMsg.Signature = hex.EncodeToString(signature.Data)
	message.Payload = consensusMsg
	return &message, nil
}

func (cp *CommitPool) IsExistingCommit(commitMsg *models.Message) bool {
	consensusMessage := commitMsg.Payload
	var exists bool
	for _, v := range cp.commitMsgs[consensusMessage.ConsensusID] {
		if v.From == commitMsg.From {
			exists = true
		}
	}
	return exists
}

func (cp *CommitPool) IsValidCommit(commit *models.Message) bool {
	consensusMsg := commit.Payload
	buf, err := hex.DecodeString(consensusMsg.Signature)
	if err != nil {
		return false
	}
	err = sigs.Verify(&types.Signature{Type: types.SigTypeEd25519, Data: buf}, commit.From, []byte(consensusMsg.Data))
	if err != nil {
		return false
	}
	return true
}

func (cp *CommitPool) AddCommit(commit *models.Message) {
	consensusID := commit.Payload.ConsensusID
	if _, ok := cp.commitMsgs[consensusID]; !ok {
		cp.commitMsgs[consensusID] = []*models.Message{}
	}

	cp.commitMsgs[consensusID] = append(cp.commitMsgs[consensusID], commit)
}

func (cp *CommitPool) CommitSize(consensusID string) int {
	if v, ok := cp.commitMsgs[consensusID]; ok {
		return len(v)
	}
	return 0
}
