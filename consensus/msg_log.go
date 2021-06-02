package consensus

import (
	types2 "github.com/Secured-Finance/dione/consensus/types"
	mapset "github.com/Secured-Finance/golang-set"
)

type MessageLog struct {
	messages          mapset.Set
	maxLogSize        int
	validationFuncMap map[types2.MessageType]func(message types2.ConsensusMessage)
}

func NewMessageLog() *MessageLog {
	msgLog := &MessageLog{
		messages:   mapset.NewSet(),
		maxLogSize: 0, // TODO
	}

	return msgLog
}

func (ml *MessageLog) AddMessage(msg types2.ConsensusMessage) {
	ml.messages.Add(msg)
}

func (ml *MessageLog) Exists(msg types2.ConsensusMessage) bool {
	return ml.messages.Contains(msg)
}

func (ml *MessageLog) Get(typ types2.MessageType, consensusID string) []*types2.ConsensusMessage {
	var result []*types2.ConsensusMessage

	for v := range ml.messages.Iter() {
		msg := v.(types2.ConsensusMessage)
		if msg.Type == typ && msg.Task.ConsensusID == consensusID {
			result = append(result, &msg)
		}
	}

	return result
}
