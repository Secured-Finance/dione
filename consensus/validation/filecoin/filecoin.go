package filecoin

import (
	"github.com/Secured-Finance/dione/types"

	"github.com/Secured-Finance/dione/consensus/validation"
	rtypes "github.com/Secured-Finance/dione/rpc/types"
)

func ValidateGetTransaction(task *types.DioneTask) error {
	//var msg ftypes.SignedMessage
	//if err := msg.UnmarshalCBOR(bytes.NewReader(payload)); err != nil {
	//	if err := msg.Message.UnmarshalCBOR(bytes.NewReader(payload)); err != nil {
	//		return xerrors.Errorf("cannot unmarshal payload: %s", err.Error())
	//	}
	//}
	//
	//if msg.Type == ftypes.MessageTypeSecp256k1 {
	//	if err := sigs.Verify(&msg.Signature, msg.Message.From.Bytes(), msg.Message.Cid().Bytes()); err != nil {
	//		logrus.Errorf("Couldn't verify transaction %v", err)
	//		return xerrors.Errorf("Couldn't verify transaction: %v")
	//	}
	//	return nil
	//} else {
	//	// TODO: BLS Signature verification
	//	return nil
	//}
	return validation.VerifyExactMatching(task)
}

func init() {
	validation.RegisterValidation(rtypes.RPCTypeFilecoin, map[string]func(*types.DioneTask) error{
		"getTransaction": ValidateGetTransaction,
	})
}
