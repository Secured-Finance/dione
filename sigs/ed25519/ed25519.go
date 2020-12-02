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
	var privKey ed25519.PrivateKey = priv
	pubKey := privKey.Public().(ed25519.PublicKey)
	return pubKey, nil
}

func (ed25519Signer) Sign(p []byte, msg []byte) ([]byte, error) {
	var privKey ed25519.PrivateKey = p
	return ed25519.Sign(privKey, msg), nil
}

func (ed25519Signer) Verify(sig []byte, a []byte, msg []byte) error {
	id, err := peer.IDFromBytes(a)
	if err != nil {
		return err
	}
	pubKey, err := id.ExtractPublicKey()
	if err != nil {
		return err
	}

	pKeyRaw, err := pubKey.Raw()
	if err != nil {
		return err
	}

	if valid := ed25519.Verify(pKeyRaw, msg, sig); !valid {
		return xerrors.Errorf("failed to verify signature")
	}
	return nil
}

func init() {
	sigs.RegisterSignature(types.SigTypeEd25519, ed25519Signer{})
}
