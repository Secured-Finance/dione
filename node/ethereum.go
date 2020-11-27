package node

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (n *Node) subscribeOnEthContracts(ctx context.Context) {
	eventChan, subscription, err := n.Ethereum.SubscribeOnOracleEvents(ctx)
	if err != nil {
		logrus.Fatal("Couldn't subscribe on ethereum contracts, exiting... ", err)
	}

	go func() {
	EventLoop:
		for {
			select {
			case event := <-eventChan:
				{
					task, err := n.Miner.MineTask(ctx, event)
					if err != nil {
						logrus.Fatal("Failed to mine task, exiting... ", err)
					}
					if task == nil {
						continue
					}
					logrus.Infof("Started new consensus round with ID: %s", event.RequestID.String())
					err = n.ConsensusManager.Propose(event.RequestID.String(), *task, event.RequestID, event.CallbackAddress)
					if err != nil {
						logrus.Errorf("Failed to propose task: %w", err)
					}
				}
			case <-ctx.Done():
				break EventLoop
			case <-subscription.Err():
				logrus.Fatal("Error with ethereum subscription, exiting... ", err)
			}
		}
	}()
}
