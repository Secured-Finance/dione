package utils

import (
	"fmt"

	"github.com/Secured-Finance/dione/blockchain/types"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

func VerifyTx(blockHeader *types.BlockHeader, tx *types.Transaction) error {
	if tx.MerkleProof == nil {
		return fmt.Errorf("block transaction doesn't have merkle proof")
	}
	txProofVerified, err := merkletree.VerifyProofUsing(tx.Hash, false, tx.MerkleProof, [][]byte{blockHeader.Hash}, keccak256.New())
	if err != nil {
		return fmt.Errorf("failed to verify tx hash merkle proof: %s", err.Error())
	}
	if !txProofVerified {
		return fmt.Errorf("transaction doesn't present in block hash merkle tree")
	}
	if !tx.ValidateHash() {
		return fmt.Errorf("transaction hash is invalid")
	}

	return nil
}
