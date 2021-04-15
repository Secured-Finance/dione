package consensus

import (
	"context"
	"sync"

	big2 "github.com/filecoin-project/go-state-types/big"

	"github.com/Secured-Finance/dione/sigs"

	"github.com/Secured-Finance/dione/rpc"

	"github.com/Secured-Finance/dione/beacon"
	"github.com/Secured-Finance/dione/contracts/dioneOracle"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/ethclient"
	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type Miner struct {
	address      peer.ID
	ethAddress   common.Address
	mutex        sync.Mutex
	beacon       beacon.BeaconNetworks
	ethClient    *ethclient.EthereumClient
	minerStake   types.BigInt
	networkStake types.BigInt
	privateKey   []byte
}

func NewMiner(
	address peer.ID,
	ethAddress common.Address,
	beacon beacon.BeaconNetworks,
	ethClient *ethclient.EthereumClient,
	privateKey []byte,
) *Miner {
	return &Miner{
		address:    address,
		ethAddress: ethAddress,
		beacon:     beacon,
		ethClient:  ethClient,
		privateKey: privateKey,
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

	m.minerStake = *mStake
	m.networkStake = *nStake

	return nil
}

func (m *Miner) GetStakeInfo(miner common.Address) (*types.BigInt, *types.BigInt, error) {
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

func (m *Miner) MineTask(ctx context.Context, event *dioneOracle.DioneOracleNewOracleRequest) (*types.DioneTask, error) {
	beaconValues, err := beacon.BeaconEntriesForTask(ctx, m.beacon)
	if err != nil {
		return nil, xerrors.Errorf("failed to get beacon entries: %w", err)
	}
	logrus.Debug("attempting to mine the task at epoch: ", beaconValues[1].Round)

	randomBase := beaconValues[1]

	if err := m.UpdateCurrentStakeInfo(); err != nil {
		return nil, xerrors.Errorf("failed to update miner stake: %w", err)
	}

	ticket, err := m.computeTicket(&randomBase)
	if err != nil {
		return nil, xerrors.Errorf("scratching ticket failed: %w", err)
	}

	winner, err := IsRoundWinner(
		types.DrandRound(randomBase.Round),
		m.address,
		randomBase,
		m.minerStake,
		m.networkStake,
		func(id peer.ID, bytes []byte) (*types.Signature, error) {
			return sigs.Sign(types.SigTypeEd25519, m.privateKey, bytes)
		},
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to check if we win next round: %w", err)
	}

	if winner == nil {
		return nil, nil
	}

	rpcMethod := rpc.GetRPCMethod(event.OriginChain, event.RequestType)
	if rpcMethod == nil {
		return nil, xerrors.Errorf("invalid rpc method name/type")
	}
	res, err := rpcMethod(event.RequestParams)
	if err != nil {
		return nil, xerrors.Errorf("couldn't do rpc request: %w", err)
	}

	return &types.DioneTask{
		OriginChain:      event.OriginChain,
		RequestType:      event.RequestType,
		RequestParams:    event.RequestParams,
		RequestID:        event.ReqID.String(),
		ConsensusID:      event.ReqID.String(),
		CallbackAddress:  event.CallbackAddress.Bytes(),
		CallbackMethodID: event.CallbackMethodID[:],
		Miner:            m.address,
		MinerEth:         m.ethAddress.Hex(),
		Ticket:           ticket,
		ElectionProof:    winner,
		BeaconEntries:    beaconValues,
		Payload:          res,
		DrandRound:       types.DrandRound(randomBase.Round),
	}, nil
}

func (m *Miner) computeTicket(brand *types.BeaconEntry) (*types.Ticket, error) {
	buf, err := m.address.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal address: %w", err)
	}

	round := types.DrandRound(brand.Round)

	input, err := DrawRandomness(brand.Data, crypto.DomainSeparationTag_TicketProduction, round-types.TicketRandomnessLookback, buf)
	if err != nil {
		return nil, err
	}

	vrfOut, err := ComputeVRF(func(id peer.ID, bytes []byte) (*types.Signature, error) {
		return sigs.Sign(types.SigTypeEd25519, m.privateKey, bytes)
	}, m.address, input)
	if err != nil {
		return nil, err
	}

	return &types.Ticket{
		VRFProof: vrfOut,
	}, nil
}

func (m *Miner) IsMinerEligibleToProposeTask(ethAddress common.Address) error {
	mStake, err := m.ethClient.GetMinerStake(ethAddress)
	if err != nil {
		return err
	}
	ok := mStake.GreaterThanEqual(big2.NewInt(ethclient.MinMinerStake))
	if !ok {
		return xerrors.Errorf("miner doesn't have enough staked tokens")
	}
	return nil
}
