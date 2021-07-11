package consensus

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/Secured-Finance/dione/beacon"

	types3 "github.com/Secured-Finance/dione/blockchain/types"

	"github.com/Secured-Finance/dione/blockchain"
	"github.com/Secured-Finance/dione/blockchain/utils"
	types2 "github.com/Secured-Finance/dione/consensus/types"
	"github.com/Secured-Finance/dione/consensus/validation"
	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/fxamacker/cbor/v2"
	"github.com/sirupsen/logrus"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

type ConsensusValidator struct {
	validationFuncMap map[types2.ConsensusMessageType]func(msg types2.ConsensusMessage, metadata map[string]interface{}) bool
	miner             *Miner
	beacon            beacon.BeaconNetwork
	blockchain        *blockchain.BlockChain
}

func NewConsensusValidator(miner *Miner, bc *blockchain.BlockChain, b beacon.BeaconNetwork) *ConsensusValidator {
	cv := &ConsensusValidator{
		miner:      miner,
		blockchain: bc,
		beacon:     b,
	}

	cv.validationFuncMap = map[types2.ConsensusMessageType]func(msg types2.ConsensusMessage, metadata map[string]interface{}) bool{
		// FIXME it all
		types2.ConsensusMessageTypePrePrepare: func(msg types2.ConsensusMessage, metadata map[string]interface{}) bool {
			// === verify block signature ===
			pubkey, err := msg.Block.Header.Proposer.ExtractPublicKey()
			if err != nil {
				logrus.Errorf("unable to extract public key from block proposer's peer id: %s", err.Error())
				return false
			}

			ok, err := pubkey.Verify(msg.Block.Header.Hash, msg.Block.Header.Signature)
			if err != nil {
				logrus.Errorf("failed to verify block signature: %s", err.Error())
				return false
			}
			if !ok {
				logrus.Errorf("signature of block %x is invalid", msg.Block.Header.Hash)
				return false
			}
			/////////////////////////////////

			// === check last hash merkle proof ===
			latestHeight, err := cv.blockchain.GetLatestBlockHeight()
			if err != nil {
				logrus.Error(err)
				return false
			}
			previousBlockHeader, err := cv.blockchain.FetchBlockHeaderByHeight(latestHeight)
			if err != nil {
				logrus.Error(err)
				return false
			}
			if bytes.Compare(msg.Block.Header.LastHash, previousBlockHeader.Hash) != 0 {
				logrus.Errorf("block header has invalid last block hash (expected: %x, actual %x)", previousBlockHeader.Hash, msg.Block.Header.LastHash)
				return false
			}

			verified, err := merkletree.VerifyProofUsing(previousBlockHeader.Hash, false, msg.Block.Header.LastHashProof, [][]byte{msg.Block.Header.Hash}, keccak256.New())
			if err != nil {
				logrus.Error("failed to verify last block hash merkle proof: %s", err.Error())
				return false
			}
			if !verified {
				logrus.Error("merkle hash of current block doesn't contain hash of previous block: %s", err.Error())
				return false
			}
			/////////////////////////////////

			// === verify election proof wincount preliminarily ===
			if msg.Block.Header.ElectionProof.WinCount < 1 {
				logrus.Error("miner isn't a winner!")
				return false
			}
			/////////////////////////////////

			// === verify miner's eligibility to propose this task ===
			err = cv.miner.IsMinerEligibleToProposeBlock(msg.Block.Header.ProposerEth)
			if err != nil {
				logrus.Errorf("miner is not eligible to propose block: %v", err)
				return false
			}
			/////////////////////////////////

			// === verify election proof vrf ===
			proposerBuf, err := msg.Block.Header.Proposer.MarshalBinary()
			if err != nil {
				logrus.Error(err)
				return false
			}

			res, err := b.Beacon.Entry(context.TODO(), msg.Block.Header.ElectionProof.RandomnessRound)
			if err != nil {
				logrus.Error(err)
				return false
			}
			eproofRandomness, err := DrawRandomness(
				res.Data,
				crypto.DomainSeparationTag_ElectionProofProduction,
				msg.Block.Header.Height,
				proposerBuf,
			)
			if err != nil {
				logrus.Errorf("failed to draw ElectionProof randomness: %s", err.Error())
				return false
			}
			err = VerifyVRF(*msg.Block.Header.Proposer, eproofRandomness, msg.Block.Header.ElectionProof.VRFProof)
			if err != nil {
				logrus.Errorf("failed to verify election proof vrf: %v", err)
				return false
			}
			//////////////////////////////////////

			// === compute wincount locally and verify values ===
			mStake, nStake, err := cv.miner.GetStakeInfo(msg.Block.Header.ProposerEth)
			if err != nil {
				logrus.Errorf("failed to get miner stake: %v", err)
				return false
			}
			actualWinCount := msg.Block.Header.ElectionProof.ComputeWinCount(mStake, nStake)
			if msg.Block.Header.ElectionProof.WinCount != actualWinCount {
				logrus.Errorf("locally computed wincount of block %x isn't matching received value!", msg.Block.Header.Hash)
				return false
			}
			//////////////////////////////////////

			// === validate block transactions ===
			result := make(chan error)
			var wg sync.WaitGroup
			for _, v := range msg.Block.Data {
				wg.Add(1)
				go func(v *types3.Transaction, c chan error) {
					defer wg.Done()
					if err := utils.VerifyTx(msg.Block.Header, v); err != nil {
						c <- fmt.Errorf("failed to verify tx: %w", err)
						return
					}

					var task types.DioneTask
					err = cbor.Unmarshal(v.Data, &task)
					if err != nil {
						c <- fmt.Errorf("failed to unmarshal transaction payload: %w", err)
						return
					}

					if validationFunc := validation.GetValidationMethod(task.OriginChain, task.RequestType); validationFunc != nil {
						if err := validationFunc(&task); err != nil {
							c <- fmt.Errorf("payload validation has been failed: %w", err)
							return
						}
					} else {
						logrus.Debugf("Origin chain [%v]/request type[%v] doesn't have any payload validation!", task.OriginChain, task.RequestType)
					}
				}(v, result)
			}
			go func() {
				wg.Wait()
				close(result)
			}()
			for err := range result {
				if err != nil {
					logrus.Error(err)
					return false
				}
			}
			/////////////////////////////////

			return true
		},
		types2.ConsensusMessageTypePrepare: func(msg types2.ConsensusMessage, metadata map[string]interface{}) bool {
			pubKey, err := msg.From.ExtractPublicKey()
			if err != nil {
				// TODO logging
				return false
			}
			ok, err := pubKey.Verify(msg.Blockhash, msg.Signature)
			if err != nil {
				// TODO logging
				return false
			}
			return ok
		},
		types2.ConsensusMessageTypeCommit: func(msg types2.ConsensusMessage, metadata map[string]interface{}) bool {
			pubKey, err := msg.From.ExtractPublicKey()
			if err != nil {
				// TODO logging
				return false
			}
			ok, err := pubKey.Verify(msg.Blockhash, msg.Signature)
			if err != nil {
				// TODO logging
				return false
			}
			return ok
		},
	}

	return cv
}

func (cv *ConsensusValidator) Valid(msg types2.ConsensusMessage, metadata map[string]interface{}) bool {
	return cv.validationFuncMap[msg.Type](msg, metadata)
}
