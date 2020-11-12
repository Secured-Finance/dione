package wallet

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
)

type Key struct {
	types.KeyInfo

	PublicKey []byte
	Address   peer.ID
}

func GenerateKey(typ types.KeyType) (*Key, error) {
	ctyp := ActSigType(typ)
	if ctyp == types.SigTypeUnknown {
		return nil, fmt.Errorf("unknown sig type: %s", typ)
	}
	pk, err := sigs.Generate(ctyp)
	if err != nil {
		return nil, err
	}
	ki := types.KeyInfo{
		Type:       typ,
		PrivateKey: pk,
	}
	return NewKey(ki)
}

// NewKey generates a new Key based on private key and signature type
// it works with both types of signatures Secp256k1 and BLS
func NewKey(keyinfo types.KeyInfo) (*Key, error) {
	k := &Key{
		KeyInfo: keyinfo,
	}

	var err error
	k.PublicKey, err = sigs.ToPublic(ActSigType(k.Type), k.PrivateKey)
	if err != nil {
		return nil, err
	}

	switch k.Type {
	case types.KTEd25519:
		pubKey, err := crypto.UnmarshalEd25519PublicKey(k.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal ed25519 public key: %w", err)
		}
		k.Address, err = peer.IDFromPublicKey(pubKey)
		if err != nil {
			return nil, fmt.Errorf("converting Secp256k1 to address: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %s", k.Type)
	}
	return k, nil

}

func ActSigType(typ types.KeyType) types.SigType {
	switch typ {
	case types.KTEd25519:
		return types.SigTypeEd25519
	default:
		return types.SigTypeUnknown
	}
}
