package consensus

import (
	"context"
	"encoding/hex"
	"time"

	"math/big"

	"github.com/Secured-Finance/dione/contracts/dioneDispute"
	"github.com/Secured-Finance/dione/contracts/dioneOracle"
	"github.com/Secured-Finance/dione/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

type DisputeManager struct {
	ctx           context.Context
	ethClient     *ethclient.EthereumClient
	pcm           *PBFTConsensusManager
	submissionMap map[string]*dioneOracle.DioneOracleSubmittedOracleRequest
	disputeMap    map[string]*dioneDispute.DioneDisputeNewDispute
	voteWindow    time.Duration
}

func NewDisputeManager(ctx context.Context, ethClient *ethclient.EthereumClient, pcm *PBFTConsensusManager, voteWindow int) (*DisputeManager, error) {
	newSubmittionsChan, submSubscription, err := ethClient.SubscribeOnNewSubmittions(ctx)
	if err != nil {
		return nil, err
	}

	newDisputesChan, dispSubscription, err := ethClient.SubscribeOnNewDisputes(ctx)
	if err != nil {
		return nil, err
	}

	dm := &DisputeManager{
		ethClient:     ethClient,
		pcm:           pcm,
		ctx:           ctx,
		submissionMap: map[string]*dioneOracle.DioneOracleSubmittedOracleRequest{},
		disputeMap:    map[string]*dioneDispute.DioneDisputeNewDispute{},
		voteWindow:    time.Duration(voteWindow) * time.Second,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					submSubscription.Unsubscribe()
					dispSubscription.Unsubscribe()
					return
				}
			case s := <-newSubmittionsChan:
				{
					dm.onNewSubmission(s)
				}
			case d := <-newDisputesChan:
				{
					dm.onNewDispute(d)
				}
			}
		}
	}()

	return dm, nil
}

func (dm *DisputeManager) onNewSubmission(submittion *dioneOracle.DioneOracleSubmittedOracleRequest) {
	c := dm.pcm.GetConsensusInfo(submittion.ReqID.String())
	if c == nil {
		// todo: warn
		return
	}

	dm.submissionMap[submittion.ReqID.String()] = submittion

	submHashBytes := sha3.Sum256(submittion.Data)
	localHashBytes := sha3.Sum256(c.Task.Payload)
	submHash := hex.EncodeToString(submHashBytes[:])
	localHash := hex.EncodeToString(localHashBytes[:])
	if submHash != localHash {
		logrus.Debugf("submission of request id %s isn't valid - beginning dispute", c.Task.RequestID)
		addr := common.HexToAddress(c.Task.MinerEth)
		reqID, ok := big.NewInt(0).SetString(c.Task.RequestID, 10)
		if !ok {
			logrus.Errorf("cannot parse request id: %s", c.Task.RequestID)
			return
		}
		err := dm.ethClient.BeginDispute(addr, reqID)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		disputeFinishTimer := time.NewTimer(dm.voteWindow)
		go func() {
			for {
				select {
				case <-dm.ctx.Done():
					return
				case <-disputeFinishTimer.C:
					{
						d, ok := dm.disputeMap[reqID.String()]
						if !ok {
							logrus.Error("cannot finish dispute: it doesn't exist in manager's dispute map!")
							return
						}
						err := dm.ethClient.FinishDispute(d.Dhash)
						if err != nil {
							logrus.Errorf(err.Error())
							return
						}
					}
				}
			}
		}()
	}
}

func (dm *DisputeManager) onNewDispute(dispute *dioneDispute.DioneDisputeNewDispute) {
	c := dm.pcm.GetConsensusInfo(dispute.RequestID.String())
	if c == nil {
		// todo: warn
		return
	}

	subm, ok := dm.submissionMap[dispute.RequestID.String()]
	if !ok {
		// todo: warn
		return
	}

	dm.disputeMap[dispute.RequestID.String()] = dispute

	if dispute.DisputeInitiator.Hex() == dm.ethClient.GetEthAddress().Hex() {
		return
	}

	submHashBytes := sha3.Sum256(subm.Data)
	localHashBytes := sha3.Sum256(c.Task.Payload)
	submHash := hex.EncodeToString(submHashBytes[:])
	localHash := hex.EncodeToString(localHashBytes[:])
	if submHash == localHash {
		err := dm.ethClient.VoteDispute(dispute.Dhash, false)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
	}

	err := dm.ethClient.VoteDispute(dispute.Dhash, true)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
}
