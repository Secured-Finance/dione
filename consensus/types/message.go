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
	_               struct{} `cbor:",toarray" hash:"-"`
	ConsensusID     string
	Signature       []byte `hash:"-"`
	RequestID       string
	CallbackAddress string
	Task            types.DioneTask
}

type Message struct {
	Type    MessageType      `cbor:"1,keyasint"`
	Payload ConsensusMessage `cbor:"2,keyasint"`
	From    peer.ID          `cbor:"-"`
}
