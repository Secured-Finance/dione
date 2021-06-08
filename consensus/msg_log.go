package consensus

import (
	"bytes"

	types2 "github.com/Secured-Finance/dione/consensus/types"
	mapset "github.com/Secured-Finance/golang-set"
)

type ConsensusMessageLog struct {
	messages   mapset.Set
	maxLogSize int
}

func NewConsensusMessageLog() *ConsensusMessageLog {
	msgLog := &ConsensusMessageLog{
		messages:   mapset.NewSet(),
		maxLogSize: 0, // TODO
	}

	return msgLog
}

func (ml *ConsensusMessageLog) AddMessage(msg types2.ConsensusMessage) {
	ml.messages.Add(msg)
}

func (ml *ConsensusMessageLog) Exists(msg types2.ConsensusMessage) bool {
	return ml.messages.Contains(msg)
}

func (ml *ConsensusMessageLog) Get(typ types2.ConsensusMessageType, blockhash []byte) []*types2.ConsensusMessage {
	var result []*types2.ConsensusMessage

	for v := range ml.messages.Iter() {
		msg := v.(types2.ConsensusMessage)
		if msg.Block != nil {

		}
		if msg.Type == typ {
			var msgBlockHash []byte
			if msg.Block != nil {
				msgBlockHash = msg.Block.Header.Hash
			} else {
				msgBlockHash = msg.Blockhash
			}
			if bytes.Compare(msgBlockHash, blockhash) == 0 {
				result = append(result, &msg)
			}
		}
	}

	return result
}
