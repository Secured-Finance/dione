package sigs

import (
	"context"
	fmt "fmt"

	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-address"
	"golang.org/x/xerrors"
)

// Sign takes in signature type, private key and message. Returns a signature for that message.
// Valid sigTypes are: "ed25519"
func Sign(sigType types.SigType, privkey []byte, msg []byte) (*types.Signature, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot sign message with signature of unsupported type: %v", sigType)
	}

	sb, err := sv.Sign(privkey, msg)
	if err != nil {
		return nil, err
	}
	return &types.Signature{
		Type: sigType,
		Data: sb,
	}, nil
}

// Verify verifies signatures
func Verify(sig *types.Signature, addrByte []byte, msg []byte) error {
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}

	sv, ok := sigs[sig.Type]
	if !ok {
		return fmt.Errorf("cannot verify signature of unsupported type: %v", sig.Type)
	}

	return sv.Verify(sig.Data, addrByte, msg)
}

// Generate generates private key of given type
func Generate(sigType types.SigType) ([]byte, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot generate private key of unsupported type: %v", sigType)
	}

	return sv.GenPrivate()
}

// ToPublic converts private key to public key
func ToPublic(sigType types.SigType, pk []byte) ([]byte, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot generate public key of unsupported type: %v", sigType)
	}

	return sv.ToPublic(pk)
}

func CheckTaskSignature(ctx context.Context, task *types.DioneTask, worker address.Address) error {
	//if task.IsValidated() {
	//	return nil
	//}
	//
	//if task.BlockSig == nil {
	//	return xerrors.New("block signature not present")
	//}
	//
	//sigb, err := task.SigningBytes()
	//if err != nil {
	//	return xerrors.Errorf("failed to get block signing bytes: %w", err)
	//}
	//
	//err = Verify(task.BlockSig, worker, sigb)
	//if err == nil {
	//	task.SetValidated()
	//}

	// TODO

	return nil
}

// SigShim is used for introducing signature functions
type SigShim interface {
	GenPrivate() ([]byte, error)
	ToPublic(pk []byte) ([]byte, error)
	Sign(pk []byte, msg []byte) ([]byte, error)
	Verify(sig []byte, a []byte, msg []byte) error
}

var sigs map[types.SigType]SigShim

// RegisterSignature should be only used during init
func RegisterSignature(typ types.SigType, vs SigShim) {
	if sigs == nil {
		sigs = make(map[types.SigType]SigShim)
	}
	sigs[typ] = vs
}
