package models

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
	ConsensusID     string
	Signature       string
	RequestID       string
	CallbackAddress string
	Data            string
}

type Message struct {
	Type    MessageType      `json:"type"`
	Payload ConsensusMessage `json:"payload"`
	From    peer.ID          `json:"-"`
}
