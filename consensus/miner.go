package consensus

import (
	"bytes"
	"context"

	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type Miner struct {
	address address.Address
	api     MinerAPI
}

type MinerAPI interface {
	WalletSign(context.Context, address.Address, []byte) (*crypto.Signature, error)
	//	TODO: get miner base based on epoch;
}

type MinerBase struct {
	MinerStake      types.BigInt
	NetworkStake    types.BigInt
	WorkerKey       address.Address
	EthWallet       common.Address
	PrevBeaconEntry types.BeaconEntry
	BeaconEntries   []types.BeaconEntry
	NullRounds      types.TaskEpoch
}

type MiningBase struct {
	epoch      types.TaskEpoch
	nullRounds types.TaskEpoch
}

func NewMinerBase(minerStake, networkStake types.BigInt, minerWallet address.Address,
	minerEthWallet common.Address, prev types.BeaconEntry, entries []types.BeaconEntry) *MinerBase {
	return &MinerBase{
		MinerStake:      minerStake,
		NetworkStake:    networkStake,
		WorkerKey:       minerWallet,
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
		BeaconEntries: mb.BeaconEntries,
		// TODO: signature
		Height: round,
	}, nil
}

func (m *Miner) computeTicket(ctx context.Context, brand *types.BeaconEntry, base *MiningBase, mb *MinerBase) (*types.Ticket, error) {
	buf := new(bytes.Buffer)
	if err := m.address.MarshalCBOR(buf); err != nil {
		return nil, xerrors.Errorf("failed to marshal address to cbor: %w", err)
	}

	round := base.epoch + base.nullRounds + 1

	input, err := DrawRandomness(brand.Data, crypto.DomainSeparationTag_TicketProduction, round-types.TicketRandomnessLookback, buf.Bytes())
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
