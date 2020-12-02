package types

import (
	"fmt"
	"io"
	"math"

	"github.com/Secured-Finance/dione/types"
	ltypes "github.com/filecoin-project/lotus/chain/types"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type SignedMessage struct {
	Message   ltypes.Message
	Signature *types.Signature
	Type      MessageType
}

type MessageType byte

const (
	MessageTypeUnknown = MessageType(math.MaxUint8)

	MessageTypeBLS       = MessageType(iota)
	MessageTypeSecp256k1 = MessageType(0x2)
)

func (t MessageType) Name() (string, error) {
	switch t {
	case MessageTypeUnknown:
		return "unknown", nil
	case MessageTypeBLS:
		return "BLS", nil
	case MessageTypeSecp256k1:
		return "Secp256k1", nil
	default:
		return "", fmt.Errorf("invalid message signature type: %d", t)
	}
}

// CBOR operations from lotus

func (t *SignedMessage) UnmarshalCBOR(r io.Reader) error {
	*t = SignedMessage{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	{

		if err := t.Message.UnmarshalCBOR(br); err != nil {
			return fmt.Errorf("unmarshaling t.Message: %w", err)
		}

	}

	{

		if err := t.Signature.UnmarshalCBOR(br); err != nil {
			return fmt.Errorf("unmarshaling t.Signature: %w", err)
		}

	}
	return nil
}

var lengthBufSignedMessage = []byte{130}

func (t *SignedMessage) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufSignedMessage); err != nil {
		return err
	}

	if err := t.Message.MarshalCBOR(w); err != nil {
		return err
	}

	if err := t.Signature.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}
