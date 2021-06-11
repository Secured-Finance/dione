package consensus

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/Secured-Finance/dione/pubsub"

	types2 "github.com/Secured-Finance/dione/consensus/types"

	"github.com/minio/blake2b-simd"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/types"
	crypto2 "github.com/filecoin-project/go-state-types/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"golang.org/x/xerrors"
)

type SignFunc func(peer.ID, []byte) (*types.Signature, error)

func ComputeVRF(privKey crypto.PrivKey, sigInput []byte) ([]byte, error) {
	return privKey.Sign(sigInput)
}

func VerifyVRF(worker peer.ID, vrfBase, vrfproof []byte) error {
	pk, err := worker.ExtractPublicKey()
	if err != nil {
		return err
	}
	ok, err := pk.Verify(vrfproof, vrfBase)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("vrf was invalid")
	}

	return nil
}

func IsRoundWinner(round uint64,
	worker peer.ID, randomness []byte, minerStake, networkStake *big.Int, privKey crypto.PrivKey) (*types.ElectionProof, error) {

	buf, err := worker.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal address: %w", err)
	}

	electionRand, err := DrawRandomness(randomness, crypto2.DomainSeparationTag_ElectionProofProduction, round, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to draw randomness: %w", err)
	}

	vrfout, err := ComputeVRF(privKey, electionRand)
	if err != nil {
		return nil, fmt.Errorf("failed to compute VRF: %w", err)
	}

	ep := &types.ElectionProof{VRFProof: vrfout}
	j := ep.ComputeWinCount(minerStake, networkStake)
	ep.WinCount = j
	if j < 1 {
		return nil, nil
	}

	return ep, nil
}

func DrawRandomness(rbase []byte, pers crypto2.DomainSeparationTag, round uint64, entropy []byte) ([]byte, error) {
	h := blake2b.New256()
	if err := binary.Write(h, binary.BigEndian, int64(pers)); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %v", err)
	}
	VRFDigest := blake2b.Sum256(rbase)
	_, err := h.Write(VRFDigest[:])
	if err != nil {
		return nil, xerrors.Errorf("hashing VRFDigest: %w", err)
	}
	if err := binary.Write(h, binary.BigEndian, round); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %v", err)
	}
	_, err = h.Write(entropy)
	if err != nil {
		return nil, xerrors.Errorf("hashing entropy: %v", err)
	}

	return h.Sum(nil), nil
}

func NewMessage(cmsg types2.ConsensusMessage, typ types2.ConsensusMessageType, privKey crypto.PrivKey) (*pubsub.GenericMessage, error) {
	var message pubsub.GenericMessage
	switch typ {
	case types2.ConsensusMessageTypePrePrepare:
		{
			message.Type = pubsub.PrePrepareMessageType
			message.Payload = types2.PrePrepareMessage{
				Block: cmsg.Block,
			}
			break
		}
	case types2.ConsensusMessageTypePrepare:
		{
			message.Type = pubsub.PrepareMessageType
			pm := types2.PrepareMessage{
				Blockhash: cmsg.Blockhash,
			}
			signature, err := privKey.Sign(cmsg.Blockhash)
			if err != nil {
				return nil, fmt.Errorf("failed to create signature: %v", err)
			}
			pm.Signature = signature
			message.Payload = pm
			break
		}
	case types2.ConsensusMessageTypeCommit:
		{
			message.Type = pubsub.CommitMessageType
			pm := types2.CommitMessage{
				Blockhash: cmsg.Blockhash,
			}
			signature, err := privKey.Sign(cmsg.Blockhash)
			if err != nil {
				return nil, fmt.Errorf("failed to create signature: %v", err)
			}
			pm.Signature = signature
			message.Payload = pm
			break
		}
	}

	return &message, nil
}
