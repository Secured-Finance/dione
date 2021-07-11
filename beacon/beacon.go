package beacon

import (
	"context"
	"fmt"

	"github.com/Secured-Finance/dione/types"
)

type BeaconResult struct {
	Entry types.BeaconEntry
	Err   error
}

type BeaconNetworks []BeaconNetwork

func (bn BeaconNetworks) BeaconNetworkForRound(e uint64) BeaconAPI {
	for i := len(bn) - 1; i >= 0; i-- {
		bp := bn[i]
		if e >= bp.Start {
			return bp.Beacon
		}
	}
	return bn[0].Beacon
}

type BeaconNetwork struct {
	Start  uint64
	Beacon BeaconAPI
}

// BeaconAPI represents a system that provides randomness.
// Other components interrogate the BeaconAPI to acquire randomness that's
// valid for a specific chain epoch. Also to verify beacon entries that have
// been posted on chain.
type BeaconAPI interface {
	Entry(context.Context, uint64) (types.BeaconEntry, error)
	VerifyEntry(types.BeaconEntry, types.BeaconEntry) error
	NewEntries() <-chan types.BeaconEntry
	LatestBeaconRound() uint64
}

// ValidateBlockBeacons is a function that verifies block randomness
func (bn BeaconNetworks) ValidateBlockBeacons(beaconNetworks BeaconNetworks, curEntry, prevEntry types.BeaconEntry) error {
	defaultBeacon := beaconNetworks.BeaconNetworkForRound(0)

	if err := defaultBeacon.VerifyEntry(curEntry, prevEntry); err != nil {
		return fmt.Errorf("beacon entry was invalid: %w", err)
	}

	return nil
}
