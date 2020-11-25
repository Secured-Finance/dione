package node

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (n *Node) subscribeOnEthContracts(ctx context.Context) {
	eventChan, subscription, err := n.Ethereum.SubscribeOnOracleEvents(ctx)
	if err != nil {
		logrus.Fatal("Can't subscribe on ethereum contracts, exiting... ", err)
	}

	go func() {
	EventLoop:
		for {
			select {
			case event := <-eventChan:
				{
					task, err := n.Miner.MineTask(ctx, event, n.Wallet.WalletSign)
					if err != nil {
						logrus.Fatal("Error with mining algorithm, exiting... ", err)
					}
					if task == nil {
						continue
					}
					logrus.Info("Started new consensus round with ID: ", task.Signature)

					err = n.ConsensusManager.Propose(event.RequestID.String(), string(task.Payload), event.RequestID, event.CallbackAddress)
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
