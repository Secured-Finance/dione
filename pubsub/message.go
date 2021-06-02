package pubsub

import "github.com/libp2p/go-libp2p-core/peer"

type PubSubMessageType int

const (
	UnknownMessageType = iota
	PrePrepareMessageType
	PrepareMessageType
	CommitMessageType
	NewTxMessageType
	NewBlockMessageType
)

type PubSubMessage struct {
	Type    PubSubMessageType
	From    peer.ID `cbor:"-"`
	Payload []byte
}
