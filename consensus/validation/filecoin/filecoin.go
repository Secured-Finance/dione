package filecoin

import (
	"bytes"

	"github.com/Secured-Finance/dione/consensus/validation"
	rtypes "github.com/Secured-Finance/dione/rpc/types"

	ftypes "github.com/Secured-Finance/dione/rpc/filecoin/types"
	"github.com/Secured-Finance/dione/sigs"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

func ValidateGetTransaction(payload []byte) error {
	var msg ftypes.SignedMessage
	if err := msg.UnmarshalCBOR(bytes.NewReader(payload)); err != nil {
		if err := msg.Message.UnmarshalCBOR(bytes.NewReader(payload)); err != nil {
			return xerrors.Errorf("cannot unmarshal payload")
		}
	}

	if msg.Type == ftypes.MessageTypeSecp256k1 {
		if err := sigs.Verify(msg.Signature, msg.Message.From.Bytes(), msg.Message.Cid().Bytes()); err != nil {
			logrus.Errorf("Couldn't verify transaction %v", err)
			return xerrors.Errorf("Couldn't verify transaction: %v")
		}
		return nil
	} else {
		// TODO: BLS Signature verification
		return nil
	}
}

func init() {
	validation.RegisterValidation(rtypes.RPCTypeFilecoin, map[string]func([]byte) error{
		"getTransaction": ValidateGetTransaction,
	})
}
