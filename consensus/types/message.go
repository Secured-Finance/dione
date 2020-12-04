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
	ConsensusID      string
	Signature        []byte `hash:"-"`
	RequestID        string
	CallbackAddress  []byte
	CallbackMethodID []byte
	Task             types.DioneTask
}

type Message struct {
	Type    MessageType
	Payload ConsensusMessage
	From    peer.ID `cbor:"-"`
}
