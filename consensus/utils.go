package consensus

import (
	"fmt"

	"github.com/Secured-Finance/dione/consensus/types"
	"github.com/Secured-Finance/dione/sigs"
	types2 "github.com/Secured-Finance/dione/types"
	"github.com/mitchellh/hashstructure/v2"
)

func verifyTaskSignature(msg types.ConsensusMessage) error {
	cHash, err := hashstructure.Hash(msg, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}
	err = sigs.Verify(
		&types2.Signature{Type: types2.SigTypeEd25519, Data: msg.Signature},
		[]byte(msg.Task.Miner.String()),
		[]byte(fmt.Sprintf("%v", cHash)),
	)
	if err != nil {
		return err
	}
	return nil
}
