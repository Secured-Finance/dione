package consensus

import (
	types2 "github.com/Secured-Finance/dione/consensus/types"
)

type ConsensusValidator struct {
	validationFuncMap map[types2.ConsensusMessageType]func(msg types2.ConsensusMessage) bool
	miner             *Miner
}

func NewConsensusValidator(miner *Miner) *ConsensusValidator {
	cv := &ConsensusValidator{
		miner: miner,
	}

	cv.validationFuncMap = map[types2.ConsensusMessageType]func(msg types2.ConsensusMessage) bool{
		// FIXME it all
		//types2.ConsensusMessageTypePrePrepare: func(msg types2.PrePrepareMessage) bool {
		//	// TODO here we need to do validation of block itself
		//
		//	// === verify task signature ===
		//	err := VerifyTaskSignature(msg.Task)
		//	if err != nil {
		//		logrus.Errorf("unable to verify signature: %v", err)
		//		return false
		//	}
		//	/////////////////////////////////
		//
		//	// === verify if request exists in cache ===
		//	var requestEvent *dioneOracle.DioneOracleNewOracleRequest
		//	err = cv.cache.Get("request_"+msg.Task.RequestID, &requestEvent)
		//	if err != nil {
		//		logrus.Errorf("the request doesn't exist in the cache or has been failed to decode: %v", err)
		//		return false
		//	}
		//
		//	if requestEvent.OriginChain != msg.Task.OriginChain ||
		//		requestEvent.RequestType != msg.Task.RequestType ||
		//		requestEvent.RequestParams != msg.Task.RequestParams {
		//
		//		logrus.Errorf("the incoming task and cached request requestEvent don't match!")
		//		return false
		//	}
		//	/////////////////////////////////
		//
		//	// === verify election proof wincount preliminarily ===
		//	if msg.Task.ElectionProof.WinCount < 1 {
		//		logrus.Error("miner isn't a winner!")
		//		return false
		//	}
		//	/////////////////////////////////
		//
		//	// === verify miner's eligibility to propose this task ===
		//	err = cv.miner.IsMinerEligibleToProposeBlock(common.HexToAddress(msg.Task.MinerEth))
		//	if err != nil {
		//		logrus.Errorf("miner is not eligible to propose task: %v", err)
		//		return false
		//	}
		//	/////////////////////////////////
		//
		//	// === verify election proof vrf ===
		//	minerAddressMarshalled, err := msg.Task.Miner.MarshalBinary()
		//	if err != nil {
		//		logrus.Errorf("failed to marshal miner address: %v", err)
		//		return false
		//	}
		//	electionProofRandomness, err := DrawRandomness(
		//		msg.Task.BeaconEntries[1].Data,
		//		crypto.DomainSeparationTag_ElectionProofProduction,
		//		msg.Task.DrandRound,
		//		minerAddressMarshalled,
		//	)
		//	if err != nil {
		//		logrus.Errorf("failed to draw electionProofRandomness: %v", err)
		//		return false
		//	}
		//	err = VerifyVRF(msg.Task.Miner, electionProofRandomness, msg.Task.ElectionProof.VRFProof)
		//	if err != nil {
		//		logrus.Errorf("failed to verify election proof vrf: %v", err)
		//	}
		//	//////////////////////////////////////
		//
		//	// === compute wincount locally and verify values ===
		//	mStake, nStake, err := cv.miner.GetStakeInfo(common.HexToAddress(msg.Task.MinerEth))
		//	if err != nil {
		//		logrus.Errorf("failed to get miner stake: %v", err)
		//		return false
		//	}
		//	actualWinCount := msg.Task.ElectionProof.ComputeWinCount(*mStake, *nStake)
		//	if msg.Task.ElectionProof.WinCount != actualWinCount {
		//		logrus.Errorf("locally computed wincount isn't matching received value!", err)
		//		return false
		//	}
		//	//////////////////////////////////////
		//
		//	// === validate payload by specific-chain checks ===
		//	if validationFunc := validation.GetValidationMethod(msg.Task.OriginChain, msg.Task.RequestType); validationFunc != nil {
		//		err := validationFunc(msg.Task.Payload)
		//		if err != nil {
		//			logrus.Errorf("payload validation has failed: %v", err)
		//			return false
		//		}
		//	} else {
		//		logrus.Debugf("Origin chain [%v]/request type[%v] doesn't have any payload validation!", msg.Task.OriginChain, msg.Task.RequestType)
		//	}
		//	/////////////////////////////////
		//
		//	return true
		//},
		types2.ConsensusMessageTypePrepare: func(msg types2.ConsensusMessage) bool {
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
		types2.ConsensusMessageTypeCommit: func(msg types2.ConsensusMessage) bool {
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

func (cv *ConsensusValidator) Valid(msg types2.ConsensusMessage) bool {
	return cv.validationFuncMap[msg.Type](msg)
}
