package node

import (
	"fmt"

	"github.com/Secured-Finance/dione/beacon"
	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/drand"
)

//	NewBeaconSchedule creates a new beacon chain schedule
func (n *Node) NewBeaconSchedule() (beacon.Schedule, error) {
	schedule := beacon.Schedule{}
	bc, err := drand.NewDrandBeacon(config.ChainGenesis, config.TaskEpochInterval, n.PubSubRouter.Pubsub)
	if err != nil {
		return nil, fmt.Errorf("creating drand beacon: %w", err)
	}
	schedule = append(schedule, beacon.BeaconPoint{Start: 0, Beacon: bc})

	return schedule, nil
}
