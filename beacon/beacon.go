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

func (bn BeaconNetworks) BeaconNetworkForEpoch(e types.TaskEpoch) BeaconAPI {
	for i := len(bn) - 1; i >= 0; i-- {
		bp := bn[i]
		if e >= bp.Start {
			return bp.Beacon
		}
	}
	return bn[0].Beacon
}

type BeaconNetwork struct {
	Start  types.TaskEpoch
	Beacon BeaconAPI
}

// BeaconAPI represents a system that provides randomness.
// Other components interrogate the BeaconAPI to acquire randomness that's
// valid for a specific chain epoch. Also to verify beacon entries that have
// been posted on chain.
type BeaconAPI interface {
	Entry(context.Context, uint64) <-chan BeaconResult
	VerifyEntry(types.BeaconEntry, types.BeaconEntry) error
	MaxBeaconRoundForEpoch(types.TaskEpoch) uint64
}

func ValidateTaskBeacons(beaconNetworks BeaconNetworks, t *types.DioneTask, prevEpoch types.TaskEpoch, prevEntry types.BeaconEntry) error {
	parentBeacon := beaconNetworks.BeaconNetworkForEpoch(prevEpoch)
	currBeacon := beaconNetworks.BeaconNetworkForEpoch(t.Epoch)
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
	bNetwork := beaconNetworks.BeaconNetworkForEpoch(t.Epoch)
	maxRound := bNetwork.MaxBeaconRoundForEpoch(t.Epoch)
	if maxRound == prevEntry.Round {
		if len(t.BeaconEntries) != 0 {
			return fmt.Errorf("expected not to have any beacon entries in this task, got %d", len(t.BeaconEntries))
		}
		return nil
	}

	if len(t.BeaconEntries) == 0 {
		return fmt.Errorf("expected to have beacon entries in this task, but didn't find any")
	}

	last := t.BeaconEntries[len(t.BeaconEntries)-1]
	if last.Round != maxRound {
		return fmt.Errorf("expected final beacon entry in task to be at round %d, got %d", maxRound, last.Round)
	}

	for i, e := range t.BeaconEntries {
		if err := bNetwork.VerifyEntry(e, prevEntry); err != nil {
			return fmt.Errorf("beacon entry %d (%d - %x (%d)) was invalid: %w", i, e.Round, e.Data, len(e.Data), err)
		}
		prevEntry = e
	}

	return nil
}

func BeaconEntriesForTask(ctx context.Context, beaconNetworks BeaconNetworks, epoch types.TaskEpoch, prevEpoch types.TaskEpoch, prev types.BeaconEntry) ([]types.BeaconEntry, error) {
	prevBeacon := beaconNetworks.BeaconNetworkForEpoch(prevEpoch)
	currBeacon := beaconNetworks.BeaconNetworkForEpoch(epoch)
	if prevBeacon != currBeacon {
		// Fork logic
		round := currBeacon.MaxBeaconRoundForEpoch(epoch)
		out := make([]types.BeaconEntry, 2)
		rch := currBeacon.Entry(ctx, round-1)
		res := <-rch
		if res.Err != nil {
			return nil, fmt.Errorf("getting entry %d returned error: %w", round-1, res.Err)
		}
		out[0] = res.Entry
		rch = currBeacon.Entry(ctx, round)
		res = <-rch
		if res.Err != nil {
			return nil, fmt.Errorf("getting entry %d returned error: %w", round, res.Err)
		}
		out[1] = res.Entry
		return out, nil
	}

	beacon := beaconNetworks.BeaconNetworkForEpoch(epoch)

	start := lib.Clock.Now()

	maxRound := beacon.MaxBeaconRoundForEpoch(epoch)
	if maxRound == prev.Round {
		return nil, nil
	}

	// TODO: this is a sketchy way to handle the genesis block not having a beacon entry
	if prev.Round == 0 {
		prev.Round = maxRound - 1
	}

	cur := maxRound
	var out []types.BeaconEntry
	for cur > prev.Round {
		rch := beacon.Entry(ctx, cur)
		select {
		case resp := <-rch:
			if resp.Err != nil {
				return nil, fmt.Errorf("beacon entry request returned error: %w", resp.Err)
			}

			out = append(out, resp.Entry)
			cur = resp.Entry.Round - 1
		case <-ctx.Done():
			return nil, fmt.Errorf("context timed out waiting on beacon entry to come back for epoch %d: %w", epoch, ctx.Err())
		}
	}

	logrus.Debug("fetching beacon entries", "took", lib.Clock.Since(start), "numEntries", len(out))
	reverse(out)
	return out, nil
}

func reverse(arr []types.BeaconEntry) {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-(1+i)] = arr[len(arr)-(1+i)], arr[i]
	}
}
