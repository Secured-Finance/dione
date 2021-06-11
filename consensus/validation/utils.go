package validation

import (
	"bytes"
	"fmt"

	"github.com/Secured-Finance/dione/rpc"
	"github.com/Secured-Finance/dione/types"
)

func VerifyExactMatching(task *types.DioneTask) error {
	rpcMethod := rpc.GetRPCMethod(task.OriginChain, task.RequestType)
	if rpcMethod == nil {
		return fmt.Errorf("invalid RPC method")
	}
	res, err := rpcMethod(task.RequestParams)
	if err != nil {
		return fmt.Errorf("failed to invoke RPC method: %w", err)
	}
	if bytes.Compare(res, task.Payload) != 0 {
		return fmt.Errorf("actual rpc response doesn't match with task's payload")
	}
	return nil
}
