package ed25519

import (
	"crypto/ed25519"

	"golang.org/x/xerrors"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/sigs"
	"github.com/Secured-Finance/dione/types"
)

type ed25519Signer struct{}

func (ed25519Signer) GenPrivate() ([]byte, error) {
	_, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return privKey.Seed(), nil
}

func (ed25519Signer) ToPublic(priv []byte) ([]byte, error) {
	privKey := ed25519.NewKeyFromSeed(priv)
	pubKey := privKey.Public().(ed25519.PublicKey)
	return pubKey, nil
}

func (ed25519Signer) Sign(p []byte, msg []byte) ([]byte, error) {
	privKey := ed25519.NewKeyFromSeed(p)
	return ed25519.Sign(privKey, msg), nil
}

func (ed25519Signer) Verify(sig []byte, a peer.ID, msg []byte) error {
	pubKey, err := a.ExtractPublicKey()
	if err != nil {
		return err
	}

	if valid, err := pubKey.Verify(msg, sig); err != nil || !valid {
		return xerrors.Errorf("failed to verify signature")
	}
	return nil
}

func init() {
	sigs.RegisterSignature(types.SigTypeEd25519, ed25519Signer{})
}
