package types

import (
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
	_               struct{} `cbor:",toarray"`
	ConsensusID     string
	Signature       []byte
	RequestID       string
	CallbackAddress string
	Data            string
}

type Message struct {
	Type    MessageType      `cbor:"1,keyasint"`
	Payload ConsensusMessage `cbor:"2,keyasint"`
	From    peer.ID          `cbor:"-"`
}
