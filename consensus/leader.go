package consensus

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/lib/sigs"
	"go.opencensus.io/trace"
	"golang.org/x/xerrors"
)

type MinerWallet interface {
	WalletSign(context.Context, address.Address, []byte) (*crypto.Signature, error)
}

type MinerBase struct {
	MinerStake      types.BigInt
	NetworkStake    types.BigInt
	WorkerKey       address.Address
	PrevBeaconEntry types.BeaconEntry
	BeaconEntries   []types.BeaconEntry
}

type SignFunc func(context.Context, address.Address, []byte) (*crypto.Signature, error)

func ComputeVRF(ctx context.Context, sign SignFunc, worker address.Address, sigInput []byte) ([]byte, error) {
	sig, err := sign(ctx, worker, sigInput)
	if err != nil {
		return nil, err
	}

	if sig.Type != crypto.SigTypeBLS {
		return nil, fmt.Errorf("miner worker address was not a BLS key")
	}

	return sig.Data, nil
}

func VerifyVRF(ctx context.Context, worker address.Address, vrfBase, vrfproof []byte) error {
	_, span := trace.StartSpan(ctx, "VerifyVRF")
	defer span.End()

	sig := &crypto.Signature{
		Type: crypto.SigTypeBLS,
		Data: vrfproof,
	}

	if err := sigs.Verify(sig, worker, vrfBase); err != nil {
		return xerrors.Errorf("vrf was invalid: %w", err)
	}

	return nil
}

func IsRoundWinner(ctx context.Context, round types.TaskEpoch,
	worker address.Address, brand types.BeaconEntry, mbi *MinerBase, a MinerWallet) (*types.ElectionProof, error) {

	buf := new(bytes.Buffer)
	if err := worker.MarshalCBOR(buf); err != nil {
		return nil, xerrors.Errorf("failed to cbor marshal address: %w", err)
	}

	electionRand, err := DrawRandomness(brand.Data, crypto.DomainSeparationTag_ElectionProofProduction, round, buf.Bytes())
	if err != nil {
		return nil, xerrors.Errorf("failed to draw randomness: %w", err)
	}

	vrfout, err := ComputeVRF(ctx, a.WalletSign, mbi.WorkerKey, electionRand)
	if err != nil {
		return nil, xerrors.Errorf("failed to compute VRF: %w", err)
	}

	ep := &types.ElectionProof{VRFProof: vrfout}
	j := ep.ComputeWinCount(mbi.MinerPower, mbi.NetworkPower)
	ep.WinCount = j
	if j < 1 {
		return nil, nil
	}

	return ep, nil
}
