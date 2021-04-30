package consensus

import (
	types2 "github.com/Secured-Finance/dione/consensus/types"
	mapset "github.com/Secured-Finance/golang-set"
)

type MessageLog struct {
	messages          mapset.Set
	maxLogSize        int
	validationFuncMap map[types2.MessageType]func(message types2.Message)
}

func NewMessageLog() *MessageLog {
	msgLog := &MessageLog{
		messages:   mapset.NewSet(),
		maxLogSize: 0, // TODO
	}

	return msgLog
}

func (ml *MessageLog) AddMessage(msg types2.Message) {
	ml.messages.Add(msg)
}

func (ml *MessageLog) Exists(msg types2.Message) bool {
	return ml.messages.Contains(msg)
}

func (ml *MessageLog) GetMessagesByTypeAndConsensusID(typ types2.MessageType, consensusID string) []types2.Message {
	var result []types2.Message

	for v := range ml.messages.Iter() {
		msg := v.(types2.Message)
		if msg.Type == typ && msg.Payload.Task.ConsensusID == consensusID {
			result = append(result, msg)
		}
	}

	return result
}
