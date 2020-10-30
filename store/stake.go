package store

import (
	"math/big"

	"github.com/Secured-Finance/dione/rpcclient"
	"github.com/ethereum/go-ethereum/common"
)

// TODO: specify store for staking mechanism
type DioneStakeInfo struct {
	MinerStake *big.Int
	TotalStake *big.Int
	Ethereum   *rpcclient.EthereumClient
}

func NewDioneStakeInfo(minerStake, totalStake *big.Int, ethereumClient *rpcclient.EthereumClient) *DioneStakeInfo {
	return &DioneStakeInfo{
		MinerStake: minerStake,
		TotalStake: totalStake,
		Ethereum:   ethereumClient,
	}
}

func (d *DioneStakeInfo) UpdateMinerStake(minerAddress common.Address) error {
	minerStake, err := d.Ethereum.GetMinerStake(minerAddress)
	if err != nil {
		return err
	}

	d.MinerStake = minerStake

	return nil
}

func (d *DioneStakeInfo) UpdateTotalStake(minerAddress common.Address) error {
	totalStake, err := d.Ethereum.GetTotalStake()
	if err != nil {
		return err
	}

	d.TotalStake = totalStake

	return nil
}
