package ethclient

import (
	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	MinMinerStake = 1000
)

// Getting total stake in DioneStaking contract, this function could
// be used for storing the total stake and veryfing the stake tokens
// on new tasks
func (c *EthereumClient) GetTotalStake() (*types.BigInt, error) {
	var b types.BigInt
	totalStake, err := c.dioneStaking.TotalStake()

	if err != nil {
		return nil, err
	}

	b.Int = totalStake
	return &b, nil
}

// Getting miner stake in DioneStaking contract, this function could
// be used for storing the miner's stake and veryfing the stake tokens
// on new tasks
func (c *EthereumClient) GetMinerStake(minerAddress common.Address) (*types.BigInt, error) {
	var b types.BigInt
	minerStake, err := c.dioneStaking.MinerStake(minerAddress)

	if err != nil {
		return nil, err
	}

	b.Int = minerStake
	return &b, nil
}
