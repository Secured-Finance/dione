package beacon

import (
	"context"
	"fmt"

	"github.com/Secured-Finance/dione/lib"
	"github.com/sirupsen/logrus"

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
	Entry(context.Context, uint64) <-chan BeaconResult
	VerifyEntry(types.BeaconEntry, types.BeaconEntry) error
	LatestBeaconRound() uint64
}

// ValidateTaskBeacons is a function that verifies dione task randomness
func ValidateTaskBeacons(beaconNetworks BeaconNetworks, t *types.DioneTask, prevEntry types.BeaconEntry) error {
	parentBeacon := beaconNetworks.BeaconNetworkForRound(t.DrandRound - 1)
	currBeacon := beaconNetworks.BeaconNetworkForRound(t.DrandRound)
	if parentBeacon != currBeacon {
		if len(t.BeaconEntries) != 2 {
			return fmt.Errorf("expected two beacon entries at beacon fork, got %d", len(t.BeaconEntries))
		}
		err := currBeacon.VerifyEntry(t.BeaconEntries[1], t.BeaconEntries[0])
		if err != nil {
			return fmt.Errorf("beacon at fork point invalid: (%v, %v): %w",
				t.BeaconEntries[1], t.BeaconEntries[0], err)
		}
		return nil
	}

	// TODO: fork logic
	bNetwork := beaconNetworks.BeaconNetworkForRound(t.DrandRound)
	if uint64(t.DrandRound) == prevEntry.Round {
		if len(t.BeaconEntries) != 0 {
			return fmt.Errorf("expected not to have any beacon entries in this task, got %d", len(t.BeaconEntries))
		}
		return nil
	}

	if len(t.BeaconEntries) == 0 {
		return fmt.Errorf("expected to have beacon entries in this task, but didn't find any")
	}

	last := t.BeaconEntries[len(t.BeaconEntries)-1]
	if last.Round != uint64(t.DrandRound) {
		return fmt.Errorf("expected final beacon entry in task to be at round %d, got %d", uint64(t.DrandRound), last.Round)
	}

	for i, e := range t.BeaconEntries {
		if err := bNetwork.VerifyEntry(e, prevEntry); err != nil {
			return fmt.Errorf("beacon entry %d (%d - %x (%d)) was invalid: %w", i, e.Round, e.Data, len(e.Data), err)
		}
		prevEntry = e
	}

	return nil
}

func BeaconEntriesForTask(ctx context.Context, beaconNetworks BeaconNetworks) ([]types.BeaconEntry, error) {
	beacon := beaconNetworks.BeaconNetworkForRound(0)
	round := beacon.LatestBeaconRound()

	start := lib.Clock.Now()

	out := make([]types.BeaconEntry, 2)
	prevBeaconEntry := beacon.Entry(ctx, round-1)
	res := <-prevBeaconEntry
	if res.Err != nil {
		return nil, fmt.Errorf("getting entry %d returned error: %w", round-1, res.Err)
	}
	out[0] = res.Entry
	curBeaconEntry := beacon.Entry(ctx, round)
	res = <-curBeaconEntry
	if res.Err != nil {
		return nil, fmt.Errorf("getting entry %d returned error: %w", round, res.Err)
	}
	out[1] = res.Entry

	logrus.Debugf("fetching beacon entries: took %v, count of entries: %v", lib.Clock.Since(start), len(out))
	return out, nil
}
