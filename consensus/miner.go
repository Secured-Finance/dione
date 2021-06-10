package consensus

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Secured-Finance/dione/blockchain/pool"

	"github.com/libp2p/go-libp2p-core/crypto"

	types2 "github.com/Secured-Finance/dione/blockchain/types"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

type Miner struct {
	address      peer.ID
	ethAddress   common.Address
	mutex        sync.Mutex
	ethClient    *ethclient.EthereumClient
	minerStake   *big.Int
	networkStake *big.Int
	privateKey   crypto.PrivKey
	mempool      *pool.Mempool
}

func NewMiner(
	address peer.ID,
	ethAddress common.Address,
	ethClient *ethclient.EthereumClient,
	privateKey crypto.PrivKey,
	mempool *pool.Mempool,
) *Miner {
	return &Miner{
		address:    address,
		ethAddress: ethAddress,
		ethClient:  ethClient,
		privateKey: privateKey,
		mempool:    mempool,
	}
}

func (m *Miner) UpdateCurrentStakeInfo() error {
	mStake, err := m.ethClient.GetMinerStake(m.ethAddress)

	if err != nil {
		logrus.Warn("Can't get miner stake", err)
		return err
	}

	nStake, err := m.ethClient.GetTotalStake()

	if err != nil {
		logrus.Warn("Can't get miner stake", err)
		return err
	}

	m.minerStake = mStake
	m.networkStake = nStake

	return nil
}

func (m *Miner) GetStakeInfo(miner common.Address) (*big.Int, *big.Int, error) {
	mStake, err := m.ethClient.GetMinerStake(miner)

	if err != nil {
		logrus.Warn("Can't get miner stake", err)
		return nil, nil, err
	}

	nStake, err := m.ethClient.GetTotalStake()

	if err != nil {
		logrus.Warn("Can't get miner stake", err)
		return nil, nil, err
	}

	return mStake, nStake, nil
}

func (m *Miner) MineBlock(randomness []byte, drandRound uint64, lastBlockHeader *types2.BlockHeader) (*types2.Block, error) {
	logrus.Debug("attempting to mine the block at epoch: ", drandRound)

	if err := m.UpdateCurrentStakeInfo(); err != nil {
		return nil, fmt.Errorf("failed to update miner stake: %w", err)
	}

	winner, err := IsRoundWinner(
		drandRound,
		m.address,
		randomness,
		m.minerStake,
		m.networkStake,
		m.privateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check if we winned in next round: %w", err)
	}

	if winner == nil {
		return nil, nil
	}

	txs := m.mempool.GetAllTransactions()
	if txs == nil {
		return nil, fmt.Errorf("there is no txes for processing") // skip new consensus round because there is no transaction for processing
	}

	newBlock, err := types2.CreateBlock(lastBlockHeader, txs, m.ethAddress, m.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create new block: %w", err)
	}

	return newBlock, nil
}

func (m *Miner) IsMinerEligibleToProposeBlock(ethAddress common.Address) error {
	mStake, err := m.ethClient.GetMinerStake(ethAddress)
	if err != nil {
		return err
	}
	if mStake.Cmp(big.NewInt(ethclient.MinMinerStake)) == -1 {
		return errors.New("miner doesn't have enough staked tokens")
	}
	return nil
}
