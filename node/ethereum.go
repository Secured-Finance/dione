package node

import (
	"context"

	"github.com/Secured-Finance/dione/consensus"
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
					logrus.Info("BlockHash for Solana transaction: ", task.BlockHash)
					logrus.Info("Started new consensus round with ID: ", task.BlockHash)
					n.ConsensusManager.NewTestConsensus(string(task.BlockHash), task.BlockHash, func(finalData string) {
						if finalData != string(task.BlockHash) {
							logrus.Warn("Expected final data to be %s, not %s", finalData)
						}
					})

					if n.ConsensusManager.Consensuses[task.BlockHash].State == consensus.ConsensusCommitted {
						logrus.Info("Consensus ID: ", task.BlockHash, " was successfull")
						logrus.Info("Submitting on-chain result: ", task.BlockHash, "for consensus ID: ", task.BlockHash)
						if err := n.Ethereum.SubmitRequestAnswer(event.RequestID, task.BlockHash, event.CallbackAddress, event.CallbackMethodID); err != nil {
							logrus.Warn("Can't submit request to ethereum chain: ", err)
						}
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
