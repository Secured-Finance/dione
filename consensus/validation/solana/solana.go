package solana

import (
	"github.com/Secured-Finance/dione/consensus/validation"
	rtypes "github.com/Secured-Finance/dione/rpc/types"
	"github.com/Secured-Finance/dione/types"
)

func ValidateGetTransaction(task *types.DioneTask) error {
	return validation.VerifyExactMatching(task)
}

func init() {
	validation.RegisterValidation(rtypes.RPCTypeSolana, map[string]func(*types.DioneTask) error{
		"getTransaction": ValidateGetTransaction,
	})
}
