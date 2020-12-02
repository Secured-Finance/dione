package consensus

import (
	"fmt"

	"github.com/Secured-Finance/dione/sigs"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-state-types/crypto"
	"golang.org/x/xerrors"
)

type SignFunc func(peer.ID, []byte) (*types.Signature, error)

func ComputeVRF(sign SignFunc, worker peer.ID, sigInput []byte) ([]byte, error) {
	sig, err := sign(worker, sigInput)
	if err != nil {
		return nil, err
	}

	if sig.Type != types.SigTypeEd25519 {
		return nil, fmt.Errorf("miner worker address was not a Ed25519 key")
	}

	return sig.Data, nil
}

func VerifyVRF(worker peer.ID, vrfBase, vrfproof []byte) error {
	err := sigs.Verify(&types.Signature{Type: types.SigTypeEd25519, Data: vrfproof}, []byte(worker), vrfBase)
	if err != nil {
		return xerrors.Errorf("vrf was invalid: %w", err)
	}

	return nil
}

func IsRoundWinner(round types.DrandRound,
	worker peer.ID, brand types.BeaconEntry, minerStake, networkStake types.BigInt, sign SignFunc) (*types.ElectionProof, error) {

	buf, err := worker.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal address: %w", err)
	}

	electionRand, err := DrawRandomness(brand.Data, crypto.DomainSeparationTag_ElectionProofProduction, round, buf)
	if err != nil {
		return nil, xerrors.Errorf("failed to draw randomness: %w", err)
	}

	vrfout, err := ComputeVRF(sign, worker, electionRand)
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
