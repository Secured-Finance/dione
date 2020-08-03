package rpcclient

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethereumClient struct {
	httpClient *ethclient.Client
	wsClient   *ethclient.Client
	logger     log.Logger
}

type EthereumClient interface {
	Connect(context.Context, string) error
	Balance(context.Context, common.Address) (*big.Int, error)
	SubscribeOnSmartContractEvents(context.Context, string)
}

func (c *ethereumClient) Connect(ctx context.Context, url string, connectionType string) error {
	client, err := ethclient.Dial(url)
	if err != nil {
		c.logger.Fatal(err)
	}
	if connectionType == "websocket" {
		c.wsClient = client
	} else {
		c.httpClient = client
	}
	return nil
}

// Balance returns the balance of the given ethereum address.
func (c *ethereumClient) Balance(ctx context.Context, address common.Address) (*big.Int, error) {
	value, err := c.httpClient.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *ethereumClient) SubscribeOnSmartContractEvents(ctx context.Context, address string) {
	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := c.wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		c.logger.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			c.logger.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}

}
