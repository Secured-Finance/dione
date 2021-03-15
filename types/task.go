package types

import (
	"strconv"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/config"
)

// DrandRound represents the round number in DRAND
type DrandRound int64

const TicketRandomnessLookback = DrandRound(1)

func (e DrandRound) String() string {
	return strconv.FormatInt(int64(e), 10)
}

// DioneTask represents the values of task computation
type DioneTask struct {
	OriginChain      uint8
	RequestType      string
	RequestParams    string
	Miner            peer.ID
	MinerEth         string
	Ticket           *Ticket
	ElectionProof    *ElectionProof
	BeaconEntries    []BeaconEntry
	DrandRound       DrandRound
	Payload          []byte
	RequestID        string
	CallbackAddress  []byte
	CallbackMethodID []byte
	ConsensusID      string
	Signature        []byte `hash:"-"`
}

func NewDioneTask(
	originChain uint8,
	requestType string,
	requestParams string,
	miner peer.ID,
	ticket *Ticket,
	electionProof *ElectionProof,
	beacon []BeaconEntry,
	drand DrandRound,
	payload []byte,
) *DioneTask {
	return &DioneTask{
		OriginChain:   originChain,
		RequestType:   requestType,
		RequestParams: requestParams,
		Miner:         miner,
		Ticket:        ticket,
		ElectionProof: electionProof,
		BeaconEntries: beacon,
		DrandRound:    drand,
		Payload:       payload,
	}
}

var tasksPerEpoch = NewInt(config.TasksPerEpoch)

const sha256bits = 256
