package node

import (
	"fmt"

	"github.com/Secured-Finance/dione/beacon"
	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/drand"
)

// NewBeaconQueue creates a new beacon chain schedule
func (n *Node) NewBeaconQueue() (beacon.BeaconNetworks, error) {
	schedule := beacon.BeaconNetworks{}
	bc, err := drand.NewDrandBeacon(config.ChainGenesis, config.TaskEpochInterval, n.PubSubRouter.Pubsub)
	if err != nil {
		return nil, fmt.Errorf("creating drand beacon: %w", err)
	}
	schedule = append(schedule, beacon.BeaconNetwork{Start: 0, Beacon: bc})

	return schedule, nil
}
