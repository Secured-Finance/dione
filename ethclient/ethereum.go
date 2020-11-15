package ethclient

import (
	"context"
	"math/big"

	"github.com/Secured-Finance/dione/contracts/aggregator"
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
	client         *ethclient.Client
	ethAddress     *common.Address
	authTransactor *bind.TransactOpts
	oracleEmitter  *oracleEmitter.OracleEmitterSession
	aggregator     *aggregator.AggregatorSession
	dioneStaking   *stakingContract.DioneStakingSession
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
}

func NewEthereumClient() *EthereumClient {
	ethereumClient := &EthereumClient{}

	return ethereumClient
}

func (c *EthereumClient) Initialize(ctx context.Context, url, privateKey, oracleEmitterContractAddress, aggregatorContractAddress, dioneStakingAddress string) error {
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

func (c *EthereumClient) SubmitRequestAnswer(reqID *big.Int, data string, callbackAddress common.Address, callbackMethodID [4]byte) error {
	// privateKey, err := crypto.HexToECDSA(private_key)
	// if err != nil {
	// 	c.Logger.Fatal("Failed to generate private key", err)
	// }

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	c.Logger.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	// }

	// publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// c.Logger.Info(hexutil.Encode(publicKeyBytes)[4:])

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// nonce, err := c.HttpClient.PendingNonceAt(ctx, fromAddress)
	// if err != nil {
	// 	c.Logger.Fatal(err)
	// }

	// gasPrice, err := c.HttpClient.SuggestGasPrice(ctx)
	// if err != nil {
	// 	c.Logger.Fatal(err)
	// }

	_, err := c.aggregator.CollectData(reqID, data, callbackAddress)
	if err != nil {
		return err
	}

	return nil
}
