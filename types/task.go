package types

import (
	"strconv"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/config"
)

type TaskType byte

const (
	EthereumTaskType = TaskType(iota)
	FilecoinTaskType
	SolanaTaskType
)

//	TaskEpoch represents the timestamp of Task computed by the Dione miner
type TaskEpoch int64

const TicketRandomnessLookback = TaskEpoch(1)

func (e TaskEpoch) String() string {
	return strconv.FormatInt(int64(e), 10)
}

//	DioneTask represents the values of task computation
//	Miner is an address of miner node
type DioneTask struct {
	Type          TaskType
	Miner         peer.ID
	Ticket        *Ticket
	ElectionProof *ElectionProof
	BeaconEntries []BeaconEntry
	Signature     *Signature
	Epoch         TaskEpoch
	Payload       []byte
}

var tasksPerEpoch = NewInt(config.TasksPerEpoch)

const sha256bits = 256
