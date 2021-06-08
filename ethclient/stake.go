package ethclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	MinMinerStake = 1000
)

// GetTotalStake for getting total stake in DioneStaking contract
func (c *EthereumClient) GetTotalStake() (*big.Int, error) {
	totalStake, err := c.dioneStaking.TotalStake()
	if err != nil {
		return nil, err
	}

	return totalStake, nil
}

// GetMinerStake for getting specified miner stake in DioneStaking contract
func (c *EthereumClient) GetMinerStake(minerAddress common.Address) (*big.Int, error) {
	minerStake, err := c.dioneStaking.MinerStake(minerAddress)
	if err != nil {
		return nil, err
	}

	return minerStake, nil
}
