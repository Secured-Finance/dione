package ethclient

import (
	"context"
	"math/big"

	"github.com/Secured-Finance/dione/contracts/dioneDispute"
	"github.com/Secured-Finance/dione/contracts/dioneOracle"
	"github.com/Secured-Finance/dione/contracts/dioneStaking"
	stakingContract "github.com/Secured-Finance/dione/contracts/dioneStaking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

//	TODO: change artifacts for other contracts
type EthereumClient struct {
	client          *ethclient.Client
	ethAddress      *common.Address
	authTransactor  *bind.TransactOpts
	dioneStaking    *stakingContract.DioneStakingSession
	disputeContract *dioneDispute.DioneDisputeSession
	dioneOracle     *dioneOracle.DioneOracleSession
}

type OracleEvent struct {
	RequestType      string
	CallbackAddress  common.Address
	CallbackMethodID [4]byte
	RequestID        *big.Int
}

type Ethereum interface {
	Initialize(ctx context.Context, url, connectionType, privateKey, oracleEmitterContractAddress, aggregatorContractAddress string) error
	Balance(context.Context, string) (*big.Int, error)
	SubmitRequestAnswer(reqID *big.Int, data string, callbackAddress common.Address, callbackMethodID [4]byte) error
	BeginDispute(miner common.Address, requestID *big.Int) error
	VoteDispute(dhash string, voteStatus bool) error
	FinishDispute(dhash string) error
	SubscribeOnNewDisputes(ctx context.Context) (chan *dioneDispute.DioneDisputeNewDispute, event.Subscription, error)
}

func NewEthereumClient() *EthereumClient {
	ethereumClient := &EthereumClient{}

	return ethereumClient
}

func (c *EthereumClient) Initialize(ctx context.Context, url, privateKey, dioneStakingAddress, disputeContractAddress, dioneOracleAddress string) error {
	client, err := ethclient.Dial(url)
	if err != nil {
		return err
	}
	c.client = client
	ecdsaKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}
	authTransactor := bind.NewKeyedTransactor(ecdsaKey)
	c.authTransactor = authTransactor
	c.ethAddress = &c.authTransactor.From

	stakingContract, err := dioneStaking.NewDioneStaking(common.HexToAddress(dioneStakingAddress), client)
	if err != nil {
		return err
	}
	oracleContract, err := dioneOracle.NewDioneOracle(common.HexToAddress(dioneOracleAddress), client)
	if err != nil {
		return err
	}
	disputeContract, err := dioneDispute.NewDioneDispute(common.HexToAddress(disputeContractAddress), client)
	if err != nil {
		return err
	}
	c.dioneStaking = &dioneStaking.DioneStakingSession{
		Contract: stakingContract,
		CallOpts: bind.CallOpts{
			Pending: true,
			From:    authTransactor.From,
			Context: context.Background(),
		},
		TransactOpts: bind.TransactOpts{
			From:     authTransactor.From,
			Signer:   authTransactor.Signer,
			GasLimit: 0,   // 0 automatically estimates gas limit
			GasPrice: nil, // nil automatically suggests gas price
			Context:  context.Background(),
		},
	}
	c.disputeContract = &dioneDispute.DioneDisputeSession{
		Contract: disputeContract,
		CallOpts: bind.CallOpts{
			Pending: true,
			From:    authTransactor.From,
			Context: context.Background(),
		},
		TransactOpts: bind.TransactOpts{
			From:     authTransactor.From,
			Signer:   authTransactor.Signer,
			GasLimit: 0,   // 0 automatically estimates gas limit
			GasPrice: nil, // nil automatically suggests gas price
			Context:  context.Background(),
		},
	}
	c.dioneOracle = &dioneOracle.DioneOracleSession{
		Contract: oracleContract,
		CallOpts: bind.CallOpts{
			Pending: true,
			From:    authTransactor.From,
			Context: context.Background(),
		},
		TransactOpts: bind.TransactOpts{
			From:     authTransactor.From,
			Signer:   authTransactor.Signer,
			GasLimit: 0,   // 0 automatically estimates gas limit
			GasPrice: nil, // nil automatically suggests gas price
			Context:  context.Background(),
		},
	}
	return nil
}

func (c *EthereumClient) GetEthAddress() *common.Address {
	return c.ethAddress
}

func (c *EthereumClient) SubscribeOnOracleEvents(ctx context.Context) (chan *dioneOracle.DioneOracleNewOracleRequest, event.Subscription, error) {
	resChan := make(chan *dioneOracle.DioneOracleNewOracleRequest)
	requestsFilter := c.dioneOracle.Contract.DioneOracleFilterer
	subscription, err := requestsFilter.WatchNewOracleRequest(&bind.WatchOpts{
		Start:   nil, //last block
		Context: ctx,
	}, resChan)
	if err != nil {
		return nil, nil, err
	}
	return resChan, subscription, err
}

func (c *EthereumClient) SubmitRequestAnswer(reqID *big.Int, callbackAddress common.Address, requestParams string, requestDeadline *big.Int, data []byte) error {
	_, err := c.dioneOracle.SubmitOracleRequest(requestParams, callbackAddress, [4]byte{}, reqID, requestDeadline, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumClient) BeginDispute(miner common.Address, requestID *big.Int) error {
	_, err := c.disputeContract.BeginDispute(miner, requestID)
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumClient) VoteDispute(dhash [32]byte, voteStatus bool) error {
	_, err := c.disputeContract.Vote(dhash, voteStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumClient) FinishDispute(dhash [32]byte) error {
	_, err := c.disputeContract.FinishDispute(dhash)
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumClient) SubscribeOnNewDisputes(ctx context.Context) (chan *dioneDispute.DioneDisputeNewDispute, event.Subscription, error) {
	resChan := make(chan *dioneDispute.DioneDisputeNewDispute)
	requestsFilter := c.disputeContract.Contract.DioneDisputeFilterer
	subscription, err := requestsFilter.WatchNewDispute(&bind.WatchOpts{
		Start:   nil, //last block
		Context: ctx,
	}, resChan, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	return resChan, subscription, err
}

func (c *EthereumClient) SubscribeOnNewSubmittions(ctx context.Context) (chan *dioneOracle.DioneOracleSubmittedOracleRequest, event.Subscription, error) {
	resChan := make(chan *dioneOracle.DioneOracleSubmittedOracleRequest)
	requestsFilter := c.dioneOracle.Contract.DioneOracleFilterer
	subscription, err := requestsFilter.WatchSubmittedOracleRequest(&bind.WatchOpts{
		Start:   nil, // last block
		Context: ctx,
	}, resChan)
	if err != nil {
		return nil, nil, err
	}
	return resChan, subscription, err
}
