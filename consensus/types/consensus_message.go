package types

import (
	types3 "github.com/Secured-Finance/dione/blockchain/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type ConsensusMessageType uint8

const (
	ConsensusMessageTypeUnknown = ConsensusMessageType(iota)

	ConsensusMessageTypePrePrepare
	ConsensusMessageTypePrepare
	ConsensusMessageTypeCommit
)

// ConsensusMessage is common struct for various consensus message types. It is stored in consensus message log.
type ConsensusMessage struct {
	Type      ConsensusMessageType
	Blockhash []byte
	Signature []byte
	Block     *types3.Block // it is optional, because not all message types have block included
	From      peer.ID
}
