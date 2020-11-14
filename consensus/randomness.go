package consensus

import (
	"encoding/binary"

	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
)

func DrawRandomness(rbase []byte, pers crypto.DomainSeparationTag, round types.DrandRound, entropy []byte) ([]byte, error) {
	h := blake2b.New256()
	if err := binary.Write(h, binary.BigEndian, int64(pers)); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	VRFDigest := blake2b.Sum256(rbase)
	_, err := h.Write(VRFDigest[:])
	if err != nil {
		return nil, xerrors.Errorf("hashing VRFDigest: %w", err)
	}
	if err := binary.Write(h, binary.BigEndian, round); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	_, err = h.Write(entropy)
	if err != nil {
		return nil, xerrors.Errorf("hashing entropy: %w", err)
	}

	return h.Sum(nil), nil
}
