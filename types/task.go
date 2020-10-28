package types

import (
	"strconv"

	"github.com/Secured-Finance/dione/config"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
)

//	TaskEpoch represents the timestamp of Task computed by the Dione miner
type TaskEpoch int64

func (e TaskEpoch) String() string {
	return strconv.FormatInt(int64(e), 10)
}

//	DioneTask represents the values of task computation
//	Miner is an address of miner node
type DioneTask struct {
	Miner         address.Address
	Ticket        *Ticket
	ElectionProof *ElectionProof
	BeaconEntries []BeaconEntry
	Signature     *crypto.Signature
	Height        TaskEpoch
}

var tasksPerEpoch = NewInt(config.TasksPerEpoch)

const sha256bits = 256
