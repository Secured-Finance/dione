package consensus

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type Miner struct {
	address peer.ID
	api     MinerAPI
	mutex   sync.Mutex
}

type MinerAPI interface {
	WalletSign(context.Context, peer.ID, []byte) (*types.Signature, error)
	//	TODO: get miner base based on epoch;
}

type MinerBase struct {
	MinerStake      types.BigInt
	NetworkStake    types.BigInt
	WorkerKey       peer.ID
	EthWallet       common.Address
	PrevBeaconEntry types.BeaconEntry
	BeaconEntries   []types.BeaconEntry
	NullRounds      types.TaskEpoch
}

type MiningBase struct {
	epoch      types.TaskEpoch
	nullRounds types.TaskEpoch // currently not used
}

func NewMinerBase(minerStake, networkStake types.BigInt, minerAddress peer.ID,
	minerEthWallet common.Address, prev types.BeaconEntry, entries []types.BeaconEntry) *MinerBase {
	return &MinerBase{
		MinerStake:      minerStake,
		NetworkStake:    networkStake,
		WorkerKey:       minerAddress,
		EthWallet:       minerEthWallet,
		PrevBeaconEntry: prev,
		BeaconEntries:   entries,
	}
}

func NewMiningBase() *MiningBase {
	return &MiningBase{
		nullRounds: 0,
	}
}

// Start, Stop mining functions

func (m *Miner) MineTask(ctx context.Context, base *MiningBase, mb *MinerBase) (*types.DioneTask, error) {
	round := base.epoch + base.nullRounds + 1
	logrus.Debug("attempting to mine the task at epoch: ", round)

	prevEntry := mb.PrevBeaconEntry

	ticket, err := m.computeTicket(ctx, &prevEntry, base, mb)
	if err != nil {
		return nil, xerrors.Errorf("scratching ticket failed: %w", err)
	}

	winner, err := IsRoundWinner(ctx, round, m.address, prevEntry, mb, m.api)
	if err != nil {
		return nil, xerrors.Errorf("failed to check if we win next round: %w", err)
	}

	if winner == nil {
		return nil, nil
	}
	return &types.DioneTask{
		Miner:         m.address,
		Ticket:        ticket,
		ElectionProof: winner,
		BeaconEntries: mb.BeaconEntries, // TODO decide what we need to do with multiple beacon entries
		// TODO: signature
		Epoch: round,
	}, nil
}

func (m *Miner) computeTicket(ctx context.Context, brand *types.BeaconEntry, base *MiningBase, mb *MinerBase) (*types.Ticket, error) {
	buf, err := m.address.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal address: %w", err)
	}

	round := base.epoch + base.nullRounds + 1

	input, err := DrawRandomness(brand.Data, crypto.DomainSeparationTag_TicketProduction, round-types.TicketRandomnessLookback, buf)
	if err != nil {
		return nil, err
	}

	vrfOut, err := ComputeVRF(ctx, m.api.WalletSign, mb.WorkerKey, input)
	if err != nil {
		return nil, err
	}

	return &types.Ticket{
		VRFProof: vrfOut,
	}, nil
}
