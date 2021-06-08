package types

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/libp2p/go-libp2p-core/peer"
)

const TicketRandomnessLookback = 1

// DioneTask represents the values of task computation
// DEPRECATED!
type DioneTask struct {
	OriginChain   uint8
	RequestType   string
	RequestParams string
	Miner         peer.ID
	MinerEth      common.Address
	ElectionProof *ElectionProof
	BeaconEntries []BeaconEntry
	DrandRound    uint64
	Payload       []byte
	RequestID     string
	ConsensusID   string
	Signature     []byte `hash:"-"`
}

func NewDioneTask(
	originChain uint8,
	requestType string,
	requestParams string,
	miner peer.ID,
	electionProof *ElectionProof,
	beacon []BeaconEntry,
	drandRound uint64,
	payload []byte,
) *DioneTask {
	return &DioneTask{
		OriginChain:   originChain,
		RequestType:   requestType,
		RequestParams: requestParams,
		Miner:         miner,
		ElectionProof: electionProof,
		BeaconEntries: beacon,
		DrandRound:    drandRound,
		Payload:       payload,
	}
}
