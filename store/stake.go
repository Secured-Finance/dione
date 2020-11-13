package store

import (
	"time"

	"github.com/Secured-Finance/dione/ethclient"
	"github.com/Secured-Finance/dione/lib"
	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation"
)

type DioneStakeInfo struct {
	ID             int
	MinerStake     *types.BigInt
	TotalStake     *types.BigInt
	MinerAddress   string
	MinerEthWallet string
	Timestamp      time.Time
	Ethereum       *ethclient.EthereumClient
}

func NewDioneStakeInfo(minerStake, totalStake *types.BigInt, minerWallet, minerEthWallet string, ethereumClient *ethclient.EthereumClient) *DioneStakeInfo {
	return &DioneStakeInfo{
		MinerStake:     minerStake,
		TotalStake:     totalStake,
		MinerAddress:   minerWallet,
		MinerEthWallet: minerEthWallet,
		Ethereum:       ethereumClient,
	}
}

func (d *DioneStakeInfo) UpdateMinerStake(minerEthAddress common.Address) error {
	minerStake, err := d.Ethereum.GetMinerStake(minerEthAddress)
	if err != nil {
		return err
	}

	d.MinerStake = minerStake

	return nil
}

func (d *DioneStakeInfo) UpdateTotalStake() error {
	totalStake, err := d.Ethereum.GetTotalStake()
	if err != nil {
		return err
	}

	d.TotalStake = totalStake

	return nil
}

// Put miner's staking information into the database
func (s *Store) CreateDioneStakeInfo(stakeStore *DioneStakeInfo) error {

	if err := stakeStore.Validate(); err != nil {
		return err
	}

	now := lib.Clock.Now()

	return s.db.QueryRow(
		"INSERT INTO staking (miner_stake, total_stake, miner_address, miner_eth_wallet, timestamp) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		stakeStore.MinerStake,
		stakeStore.TotalStake,
		stakeStore.MinerAddress,
		stakeStore.MinerEthWallet,
		now,
	).Scan(&stakeStore.ID)
}

func (s *Store) GetLastStakeInfo(wallet, ethWallet string) (*DioneStakeInfo, error) {
	var stake *DioneStakeInfo
	if err := s.db.Select(&stake,
		`SELECT miner_stake, total_stake, miner_address, miner_eth_wallet, timestamp FROM staking ORDER BY TIMESTAMP DESC LIMIT 1 WHERE miner_address=$1, miner_eth_wallet=$2`,
		wallet,
		ethWallet,
	); err != nil {
		return nil, err
	}

	return stake, nil
}

// Before puting the data into the database validating all required fields
func (s *DioneStakeInfo) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.MinerStake, validation.Required, validation.By(types.ValidateBigInt(s.MinerStake.Int))),
		validation.Field(&s.TotalStake, validation.Required, validation.By(types.ValidateBigInt(s.TotalStake.Int))),
		validation.Field(&s.MinerAddress, validation.Required),
		validation.Field(&s.MinerEthWallet, validation.Required),
	)
}
