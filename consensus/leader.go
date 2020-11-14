package consensus

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-state-types/crypto"
	"golang.org/x/xerrors"
)

type SignFunc func(context.Context, peer.ID, []byte) (*types.Signature, error)

func ComputeVRF(ctx context.Context, sign SignFunc, worker peer.ID, sigInput []byte) ([]byte, error) {
	sig, err := sign(ctx, worker, sigInput)
	if err != nil {
		return nil, err
	}

	if sig.Type != types.SigTypeEd25519 {
		return nil, fmt.Errorf("miner worker address was not a Ed25519 key")
	}

	return sig.Data, nil
}

func VerifyVRF(ctx context.Context, worker peer.ID, vrfBase, vrfproof []byte) error {
	pKey, err := worker.ExtractPublicKey()
	if err != nil {
		return xerrors.Errorf("failed to extract public key from worker address: %w", err)
	}

	valid, err := pKey.Verify(vrfBase, vrfproof)
	if err != nil || !valid {
		return xerrors.Errorf("vrf was invalid: %w", err)
	}

	return nil
}

func IsRoundWinner(ctx context.Context, round types.DrandRound,
	worker peer.ID, brand types.BeaconEntry, minerStake, networkStake types.BigInt, a MinerAPI) (*types.ElectionProof, error) {

	buf, err := worker.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal address: %w", err)
	}

	electionRand, err := DrawRandomness(brand.Data, crypto.DomainSeparationTag_ElectionProofProduction, round, buf)
	if err != nil {
		return nil, xerrors.Errorf("failed to draw randomness: %w", err)
	}

	vrfout, err := ComputeVRF(ctx, a.WalletSign, worker, electionRand)
	if err != nil {
		return nil, xerrors.Errorf("failed to compute VRF: %w", err)
	}

	ep := &types.ElectionProof{VRFProof: vrfout}
	j := ep.ComputeWinCount(minerStake, networkStake)
	ep.WinCount = j
	if j < 1 {
		return nil, nil
	}

	return ep, nil
}
