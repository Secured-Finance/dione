package consensus

import (
	"encoding/binary"
	"fmt"

	"github.com/fxamacker/cbor/v2"

	"github.com/Secured-Finance/dione/pubsub"

	types2 "github.com/Secured-Finance/dione/consensus/types"

	"github.com/mitchellh/hashstructure/v2"

	"github.com/Secured-Finance/dione/sigs"
	"github.com/minio/blake2b-simd"

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

func DrawRandomness(rbase []byte, pers crypto.DomainSeparationTag, round types.DrandRound, entropy []byte) ([]byte, error) {
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

func VerifyTaskSignature(task types.DioneTask) error {
	cHash, err := hashstructure.Hash(task, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}
	err = sigs.Verify(
		&types.Signature{Type: types.SigTypeEd25519, Data: task.Signature},
		[]byte(task.Miner),
		[]byte(fmt.Sprintf("%v", cHash)),
	)
	if err != nil {
		return err
	}
	return nil
}

func NewMessage(msg *pubsub.GenericMessage, typ pubsub.PubSubMessageType) (pubsub.GenericMessage, error) {
	var newMsg pubsub.GenericMessage
	newMsg.Type = typ
	newCMsg := msg.Payload
	newMsg.Payload = newCMsg
	return newMsg, nil
}

func CreatePrePrepareWithTaskSignature(task *types.DioneTask, privateKey []byte) (*pubsub.GenericMessage, error) {
	var message pubsub.GenericMessage
	message.Type = pubsub.PrePrepareMessageType

	cHash, err := hashstructure.Hash(task, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, err
	}
	signature, err := sigs.Sign(types.SigTypeEd25519, privateKey, []byte(fmt.Sprintf("%v", cHash)))
	if err != nil {
		return nil, err
	}
	task.Signature = signature.Data
	data, err := cbor.Marshal(types2.ConsensusMessage{Task: *task})
	if err != nil {
		return nil, err
	}
	message.Payload = data
	return &message, nil
}
