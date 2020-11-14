package node

import (
	"fmt"

	"github.com/Secured-Finance/dione/types"

	"github.com/Secured-Finance/dione/beacon"
	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/drand"
)

// NewBeaconClient creates a new beacon chain client
func (n *Node) NewBeaconClient() (beacon.BeaconNetworks, error) {
	networks := beacon.BeaconNetworks{}
	bc, err := drand.NewDrandBeacon(config.ChainGenesis, config.TaskEpochInterval, n.PubSubRouter.Pubsub)
	if err != nil {
		return nil, fmt.Errorf("creating drand beacon: %w", err)
	}
	networks = append(networks, beacon.BeaconNetwork{Start: types.DrandRound(config.ChainGenesis), Beacon: bc})
	// NOTE: currently we use only one network

	return networks, nil
}
