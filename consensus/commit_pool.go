package consensus

import (
	types2 "github.com/Secured-Finance/dione/consensus/types"
	"github.com/sirupsen/logrus"
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
	newCMsg := prepareMsg.Payload
	message.Payload = newCMsg
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
	err := verifyTaskSignature(consensusMsg)
	if err != nil {
		logrus.Errorf("failed to verify task signature: %v", err)
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
