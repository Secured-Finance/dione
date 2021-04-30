package consensus

import (
	"github.com/Secured-Finance/dione/cache"
	types2 "github.com/Secured-Finance/dione/consensus/types"
	"github.com/Secured-Finance/dione/consensus/validation"
	"github.com/Secured-Finance/dione/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/sirupsen/logrus"
)

type ConsensusValidator struct {
	validationFuncMap map[types2.MessageType]func(msg types2.Message) bool
	eventCache        cache.EventCache
	miner             *Miner
}

func NewConsensusValidator(ec cache.EventCache, miner *Miner) *ConsensusValidator {
	cv := &ConsensusValidator{
		eventCache: ec,
		miner:      miner,
	}

	cv.validationFuncMap = map[types2.MessageType]func(msg types2.Message) bool{
		types2.MessageTypePrePrepare: func(msg types2.Message) bool {
			// TODO here we need to do validation of tx itself
			consensusMsg := msg.Payload

			// === verify task signature ===
			err := VerifyTaskSignature(consensusMsg.Task)
			if err != nil {
				logrus.Errorf("unable to verify signature: %v", err)
				return false
			}
			/////////////////////////////////

			// === verify if request exists in event log cache ===
			requestEvent, err := cv.eventCache.GetOracleRequestEvent("request_" + consensusMsg.Task.RequestID)
			if err != nil {
				logrus.Errorf("the incoming request task event doesn't exist in the EVC, or is broken: %v", err)
				return false
			}
			if requestEvent.OriginChain != consensusMsg.Task.OriginChain ||
				requestEvent.RequestType != consensusMsg.Task.RequestType ||
				requestEvent.RequestParams != consensusMsg.Task.RequestParams {

				logrus.Errorf("the incoming task and cached request event don't match!")
				return false
			}
			/////////////////////////////////

			// === verify election proof wincount preliminarily ===
			if consensusMsg.Task.ElectionProof.WinCount < 1 {
				logrus.Error("miner isn't a winner!")
				return false
			}
			/////////////////////////////////

			// === verify miner's eligibility to propose this task ===
			err = cv.miner.IsMinerEligibleToProposeTask(common.HexToAddress(consensusMsg.Task.MinerEth))
			if err != nil {
				logrus.Errorf("miner is not eligible to propose task: %v", err)
				return false
			}
			/////////////////////////////////

			// === verify election proof vrf ===
			minerAddressMarshalled, err := consensusMsg.Task.Miner.MarshalBinary()
			if err != nil {
				logrus.Errorf("failed to marshal miner address: %v", err)
				return false
			}
			electionProofRandomness, err := DrawRandomness(
				consensusMsg.Task.BeaconEntries[1].Data,
				crypto.DomainSeparationTag_ElectionProofProduction,
				consensusMsg.Task.DrandRound,
				minerAddressMarshalled,
			)
			if err != nil {
				logrus.Errorf("failed to draw electionProofRandomness: %v", err)
				return false
			}
			err = VerifyVRF(consensusMsg.Task.Miner, electionProofRandomness, consensusMsg.Task.ElectionProof.VRFProof)
			if err != nil {
				logrus.Errorf("failed to verify election proof vrf: %v", err)
			}
			//////////////////////////////////////

			// === verify ticket vrf ===
			ticketRandomness, err := DrawRandomness(
				consensusMsg.Task.BeaconEntries[1].Data,
				crypto.DomainSeparationTag_TicketProduction,
				consensusMsg.Task.DrandRound-types.TicketRandomnessLookback,
				minerAddressMarshalled,
			)
			if err != nil {
				logrus.Errorf("failed to draw ticket electionProofRandomness: %v", err)
				return false
			}

			err = VerifyVRF(consensusMsg.Task.Miner, ticketRandomness, consensusMsg.Task.Ticket.VRFProof)
			if err != nil {
				logrus.Errorf("failed to verify ticket vrf: %v", err)
			}
			//////////////////////////////////////

			// === compute wincount locally and verify values ===
			mStake, nStake, err := cv.miner.GetStakeInfo(common.HexToAddress(consensusMsg.Task.MinerEth))
			if err != nil {
				logrus.Errorf("failed to get miner stake: %v", err)
				return false
			}
			actualWinCount := consensusMsg.Task.ElectionProof.ComputeWinCount(*mStake, *nStake)
			if consensusMsg.Task.ElectionProof.WinCount != actualWinCount {
				logrus.Errorf("locally computed wincount isn't matching received value!", err)
				return false
			}
			//////////////////////////////////////

			// === validate payload by specific-chain checks ===
			if validationFunc := validation.GetValidationMethod(consensusMsg.Task.OriginChain, consensusMsg.Task.RequestType); validationFunc != nil {
				err := validationFunc(consensusMsg.Task.Payload)
				if err != nil {
					logrus.Errorf("payload validation has failed: %v", err)
					return false
				}
			} else {
				logrus.Debugf("Origin chain [%v]/request type[%v] doesn't have any payload validation!", consensusMsg.Task.OriginChain, consensusMsg.Task.RequestType)
			}
			/////////////////////////////////

			return true
		},
		types2.MessageTypePrepare: func(msg types2.Message) bool {
			err := VerifyTaskSignature(msg.Payload.Task)
			if err != nil {
				return false
			}
			return true
		},
		types2.MessageTypeCommit: func(msg types2.Message) bool {
			err := VerifyTaskSignature(msg.Payload.Task)
			if err != nil {
				return false
			}
			return true
		},
	}

	return cv
}

func (cv *ConsensusValidator) Valid(msg types2.Message) bool {
	return cv.validationFuncMap[msg.Type](msg)
}
