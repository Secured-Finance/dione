package consensus

import (
	"github.com/Secured-Finance/dione/cache"
	types2 "github.com/Secured-Finance/dione/consensus/types"
	"github.com/Secured-Finance/dione/consensus/validation"
	"github.com/Secured-Finance/dione/contracts/dioneOracle"
	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/sirupsen/logrus"
)

type ConsensusValidator struct {
	validationFuncMap map[types2.MessageType]func(msg types2.ConsensusMessage) bool
	cache             cache.Cache
	miner             *Miner
}

func NewConsensusValidator(ec cache.Cache, miner *Miner) *ConsensusValidator {
	cv := &ConsensusValidator{
		cache: ec,
		miner: miner,
	}

	cv.validationFuncMap = map[types2.MessageType]func(msg types2.ConsensusMessage) bool{
		types2.MessageTypePrePrepare: func(msg types2.ConsensusMessage) bool {
			// TODO here we need to do validation of tx itself

			// === verify task signature ===
			err := VerifyTaskSignature(msg.Task)
			if err != nil {
				logrus.Errorf("unable to verify signature: %v", err)
				return false
			}
			/////////////////////////////////

			// === verify if request exists in cache ===
			var requestEvent *dioneOracle.DioneOracleNewOracleRequest
			err = cv.cache.Get("request_"+msg.Task.RequestID, &requestEvent)
			if err != nil {
				logrus.Errorf("the request doesn't exist in the cache or has been failed to decode: %v", err)
				return false
			}

			if requestEvent.OriginChain != msg.Task.OriginChain ||
				requestEvent.RequestType != msg.Task.RequestType ||
				requestEvent.RequestParams != msg.Task.RequestParams {

				logrus.Errorf("the incoming task and cached request requestEvent don't match!")
				return false
			}
			/////////////////////////////////

			// === verify election proof wincount preliminarily ===
			if msg.Task.ElectionProof.WinCount < 1 {
				logrus.Error("miner isn't a winner!")
				return false
			}
			/////////////////////////////////

			// === verify miner's eligibility to propose this task ===
			err = cv.miner.IsMinerEligibleToProposeTask(common.HexToAddress(msg.Task.MinerEth))
			if err != nil {
				logrus.Errorf("miner is not eligible to propose task: %v", err)
				return false
			}
			/////////////////////////////////

			// === verify election proof vrf ===
			minerAddressMarshalled, err := msg.Task.Miner.MarshalBinary()
			if err != nil {
				logrus.Errorf("failed to marshal miner address: %v", err)
				return false
			}
			electionProofRandomness, err := DrawRandomness(
				msg.Task.BeaconEntries[1].Data,
				crypto.DomainSeparationTag_ElectionProofProduction,
				msg.Task.DrandRound,
				minerAddressMarshalled,
			)
			if err != nil {
				logrus.Errorf("failed to draw electionProofRandomness: %v", err)
				return false
			}
			err = VerifyVRF(msg.Task.Miner, electionProofRandomness, msg.Task.ElectionProof.VRFProof)
			if err != nil {
				logrus.Errorf("failed to verify election proof vrf: %v", err)
			}
			//////////////////////////////////////

			// === verify ticket vrf ===
			ticketRandomness, err := DrawRandomness(
				msg.Task.BeaconEntries[1].Data,
				crypto.DomainSeparationTag_TicketProduction,
				msg.Task.DrandRound-types.TicketRandomnessLookback,
				minerAddressMarshalled,
			)
			if err != nil {
				logrus.Errorf("failed to draw ticket electionProofRandomness: %v", err)
				return false
			}

			err = VerifyVRF(msg.Task.Miner, ticketRandomness, msg.Task.Ticket.VRFProof)
			if err != nil {
				logrus.Errorf("failed to verify ticket vrf: %v", err)
			}
			//////////////////////////////////////

			// === compute wincount locally and verify values ===
			mStake, nStake, err := cv.miner.GetStakeInfo(common.HexToAddress(msg.Task.MinerEth))
			if err != nil {
				logrus.Errorf("failed to get miner stake: %v", err)
				return false
			}
			actualWinCount := msg.Task.ElectionProof.ComputeWinCount(*mStake, *nStake)
			if msg.Task.ElectionProof.WinCount != actualWinCount {
				logrus.Errorf("locally computed wincount isn't matching received value!", err)
				return false
			}
			//////////////////////////////////////

			// === validate payload by specific-chain checks ===
			if validationFunc := validation.GetValidationMethod(msg.Task.OriginChain, msg.Task.RequestType); validationFunc != nil {
				err := validationFunc(msg.Task.Payload)
				if err != nil {
					logrus.Errorf("payload validation has failed: %v", err)
					return false
				}
			} else {
				logrus.Debugf("Origin chain [%v]/request type[%v] doesn't have any payload validation!", msg.Task.OriginChain, msg.Task.RequestType)
			}
			/////////////////////////////////

			return true
		},
		types2.MessageTypePrepare: func(msg types2.ConsensusMessage) bool {
			err := VerifyTaskSignature(msg.Task)
			if err != nil {
				return false
			}
			return true
		},
		types2.MessageTypeCommit: func(msg types2.ConsensusMessage) bool {
			err := VerifyTaskSignature(msg.Task)
			if err != nil {
				return false
			}
			return true
		},
	}

	return cv
}

func (cv *ConsensusValidator) Valid(msg types2.ConsensusMessage) bool {
	return cv.validationFuncMap[msg.Type](msg)
}
