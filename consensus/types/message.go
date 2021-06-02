package types

import (
	"github.com/Secured-Finance/dione/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type MessageType uint8

const (
	MessageTypeUnknown = MessageType(iota)

	MessageTypePrePrepare
	MessageTypePrepare
	MessageTypeCommit
)

type ConsensusMessage struct {
	Task types.DioneTask
	From peer.ID
	Type MessageType
}
