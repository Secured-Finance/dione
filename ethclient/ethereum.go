package ethclient

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/Secured-Finance/dione/contracts/aggregator"
	"github.com/Secured-Finance/dione/contracts/dioneDispute"
	"github.com/Secured-Finance/dione/contracts/dioneStaking"
	stakingContract "github.com/Secured-Finance/dione/contracts/dioneStaking"
	oracleEmitter "github.com/Secured-Finance/dione/contracts/oracleemitter"
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
	oracleEmitter   *oracleEmitter.OracleEmitterSession
	aggregator      *aggregator.AggregatorSession
	dioneStaking    *stakingContract.DioneStakingSession
	disputeContract *dioneDispute.DioneDisputeSession
	// dioneOracle    *dioneOracle.DioneOracleSession
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
	SubscribeOnSmartContractEvents(context.Context, string)
	SubmitRequestAnswer(reqID *big.Int, data string, callbackAddress common.Address, callbackMethodID [4]byte) error
	BeginDispute(miner common.Address, requestID *big.Int) error
	VoteDispute() error
}

func NewEthereumClient() *EthereumClient {
	ethereumClient := &EthereumClient{}

	return ethereumClient
}

func (c *EthereumClient) Initialize(ctx context.Context, url, privateKey, oracleEmitterContractAddress, aggregatorContractAddress, dioneStakingAddress, disputeContractAddress string) error {
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

	emitter, err := oracleEmitter.NewOracleEmitter(common.HexToAddress(oracleEmitterContractAddress), client)
	if err != nil {
		return err
	}
	aggregatorPlainSC, err := aggregator.NewAggregator(common.HexToAddress(aggregatorContractAddress), client)
	if err != nil {
		return err
	}
	stakingContract, err := dioneStaking.NewDioneStaking(common.HexToAddress(dioneStakingAddress), client)
	if err != nil {
		return err
	}
	disputeContract, err := dioneDispute.NewDioneDispute(common.HexToAddress(disputeContractAddress), client)
	if err != nil {
		return err
	}
	// oracleContract, err := dioneOracle.NewDioneOracle(common.HexToAddress(dioneOracleContract), client)
	// if err != nil {
	// 	return err
	// }
	c.oracleEmitter = &oracleEmitter.OracleEmitterSession{
		Contract: emitter,
		CallOpts: bind.CallOpts{
			Pending: true,
			From:    authTransactor.From,
			Context: context.Background(),
		},
		TransactOpts: bind.TransactOpts{
			From:     authTransactor.From,
			Signer:   authTransactor.Signer,
			GasLimit: 200000, // 0 automatically estimates gas limit
			GasPrice: nil,    // nil automatically suggests gas price
			Context:  context.Background(),
		},
	}
	c.aggregator = &aggregator.AggregatorSession{
		Contract: aggregatorPlainSC,
		CallOpts: bind.CallOpts{
			Pending: true,
			From:    authTransactor.From,
			Context: context.Background(),
		},
		TransactOpts: bind.TransactOpts{
			From:     authTransactor.From,
			Signer:   authTransactor.Signer,
			GasLimit: 200000, // 0 automatically estimates gas limit
			GasPrice: nil,    // nil automatically suggests gas price
			Context:  context.Background(),
		},
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
			GasLimit: 200000,                 // 0 automatically estimates gas limit
			GasPrice: big.NewInt(1860127603), // nil automatically suggests gas price
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
	// c.dioneOracle = &dioneOracle.DioneOracleSession{
	// 	Contract: oracleContract,
	// 	CallOpts: bind.CallOpts{
	// 		Pending: true,
	// 		From:    authTransactor.From,
	// 		Context: context.Background(),
	// 	},
	// 	TransactOpts: bind.TransactOpts{
	// 		From:     authTransactor.From,
	// 		Signer:   authTransactor.Signer,
	// 		GasLimit: 200000,                 // 0 automatically estimates gas limit
	// 		GasPrice: big.NewInt(1860127603), // nil automatically suggests gas price
	// 		Context:  context.Background(),
	// 	},
	// }
	return nil
}

func (c *EthereumClient) GetEthAddress() *common.Address {
	return c.ethAddress
}

func (c *EthereumClient) SubscribeOnOracleEvents(ctx context.Context) (chan *oracleEmitter.OracleEmitterNewOracleRequest, event.Subscription, error) {
	resChan := make(chan *oracleEmitter.OracleEmitterNewOracleRequest)
	requestsFilter := c.oracleEmitter.Contract.OracleEmitterFilterer
	subscription, err := requestsFilter.WatchNewOracleRequest(&bind.WatchOpts{
		Start:   nil, //last block
		Context: ctx,
	}, resChan)
	if err != nil {
		return nil, nil, err
	}
	return resChan, subscription, err
}

// func (c *EthereumClient) SubscribeOnSumbittedRequests(ctx context.Context) (chan *dioneOracle.DioneOracleSubmittedOracleRequest, event.Subscription, error) {
// 	resChan := make(chan *dioneOracle.DioneOracleSubmittedOracleRequest)
// 	requestsFilter := c.dioneOracle.Contract.DioneOracleFilterer
// 	subscription, err := requestsFilter.WatchSubmittedOracleRequest(&bind.WatchOpts{
// 		Start:   nil, //last block
// 		Context: ctx,
// 	}, resChan)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return resChan, subscription, err
// }

func (c *EthereumClient) SubmitRequestAnswer(reqID *big.Int, data string, callbackAddress common.Address) error {
	_, err := c.aggregator.CollectData(reqID, data, callbackAddress)
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

func (c *EthereumClient) VoteDispute(dhash string, voteStatus bool) error {
	dhashRawSlice, err := hex.DecodeString(dhash)
	if err != nil {
		return err
	}
	var dhashRaw [32]byte
	copy(dhashRaw[:], dhashRawSlice)
	_, err = c.disputeContract.Vote(dhashRaw, voteStatus)
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumClient) FinishDispute(dhash string) error {
	dhashRawSlice, err := hex.DecodeString(dhash)
	if err != nil {
		return err
	}
	var dhashRaw [32]byte
	copy(dhashRaw[:], dhashRawSlice)
	_, err = c.disputeContract.FinishDispute(dhashRaw)
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
