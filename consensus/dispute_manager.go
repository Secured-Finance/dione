package consensus

import (
	"context"
	"encoding/hex"
	"time"

	types2 "github.com/Secured-Finance/dione/blockchain/types"

	"github.com/Secured-Finance/dione/types"
	"github.com/fxamacker/cbor/v2"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	"github.com/Secured-Finance/dione/blockchain"

	"github.com/Secured-Finance/dione/contracts/dioneDispute"
	"github.com/Secured-Finance/dione/contracts/dioneOracle"
	"github.com/Secured-Finance/dione/ethclient"
)

type DisputeManager struct {
	ctx           context.Context
	ethClient     *ethclient.EthereumClient
	pcm           *PBFTConsensusManager
	submissionMap map[string]*dioneOracle.DioneOracleSubmittedOracleRequest
	disputeMap    map[string]*dioneDispute.DioneDisputeNewDispute
	voteWindow    time.Duration
	blockchain    *blockchain.BlockChain
}

func NewDisputeManager(ctx context.Context, ethClient *ethclient.EthereumClient, pcm *PBFTConsensusManager, voteWindow int, bc *blockchain.BlockChain) (*DisputeManager, error) {
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
		blockchain:    bc,
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

func (dm *DisputeManager) onNewSubmission(submission *dioneOracle.DioneOracleSubmittedOracleRequest) {
	// find a block that contains the dione task with specified request id
	task, block, err := dm.findTaskAndBlockWithRequestID(submission.ReqID.String())
	if err != nil {
		logrus.Error(err)
		return
	}

	dm.submissionMap[submission.ReqID.String()] = submission

	submHashBytes := sha3.Sum256(submission.Data)
	localHashBytes := sha3.Sum256(task.Payload)
	submHash := hex.EncodeToString(submHashBytes[:])
	localHash := hex.EncodeToString(localHashBytes[:])
	if submHash != localHash {
		logrus.Debugf("submission of request id %s isn't valid - beginning dispute", submission.ReqID)
		err := dm.ethClient.BeginDispute(block.Header.ProposerEth, submission.ReqID)
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
						d, ok := dm.disputeMap[submission.ReqID.String()]
						if !ok {
							logrus.Error("cannot finish dispute: it doesn't exist in manager's dispute map!")
							return
						}
						err := dm.ethClient.FinishDispute(d.Dhash)
						if err != nil {
							logrus.Errorf(err.Error())
							return
						}
						disputeFinishTimer.Stop()
						return
					}
				}
			}
		}()
	}
}

func (dm *DisputeManager) findTaskAndBlockWithRequestID(requestID string) (*types.DioneTask, *types2.Block, error) {
	height, err := dm.blockchain.GetLatestBlockHeight()
	if err != nil {
		return nil, nil, err
	}

	for {
		block, err := dm.blockchain.FetchBlockByHeight(height)
		if err != nil {
			return nil, nil, err
		}

		for _, v := range block.Data {
			var task types.DioneTask
			err := cbor.Unmarshal(v.Data, &task)
			if err != nil {
				logrus.Error(err)
				continue
			}

			if task.RequestID == requestID {
				return &task, block, nil
			}
		}

		height--
	}
}

func (dm *DisputeManager) onNewDispute(dispute *dioneDispute.DioneDisputeNewDispute) {
	task, _, err := dm.findTaskAndBlockWithRequestID(dispute.RequestID.String())
	if err != nil {
		logrus.Error(err)
		return
	}

	subm, ok := dm.submissionMap[dispute.RequestID.String()]
	if !ok {
		logrus.Warn("desired submission isn't found in map")
		return
	}

	dm.disputeMap[dispute.RequestID.String()] = dispute

	if dispute.DisputeInitiator.Hex() == dm.ethClient.GetEthAddress().Hex() {
		return
	}

	submHashBytes := sha3.Sum256(subm.Data)
	localHashBytes := sha3.Sum256(task.Payload)
	submHash := hex.EncodeToString(submHashBytes[:])
	localHash := hex.EncodeToString(localHashBytes[:])
	if submHash == localHash {
		err := dm.ethClient.VoteDispute(dispute.Dhash, false)
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
	}

	err = dm.ethClient.VoteDispute(dispute.Dhash, true)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
}
