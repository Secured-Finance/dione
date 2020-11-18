package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumRPCClient struct {
	client *ethclient.Client
}

func NewEthereumRPCClient(url string) (*EthereumRPCClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthereumRPCClient{
		client: client,
	}, nil
}

func (erc *EthereumRPCClient) GetTransaction(txHash string) ([]byte, error) {
	txHHash := common.HexToHash(txHash)
	tx, _, err := erc.client.TransactionByHash(context.TODO(), txHHash)
	if err != nil {
		return nil, err
	}
	txRaw, err := tx.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return txRaw, nil
}
