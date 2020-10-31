package rpcclient

import (
	"context"
	"math/big"

	stakingContract "github.com/Secured-Finance/dione/contracts/DioneStaking"
	"github.com/Secured-Finance/dione/contracts/aggregator"
	"github.com/Secured-Finance/dione/contracts/oracleemitter"
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
	oracleEmitter  *oracleemitter.OracleEmitterSession
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

func (c *EthereumClient) Initialize(ctx context.Context, url, privateKey, oracleEmitterContractAddress, aggregatorContractAddress string) error {
	client, err := ethclient.Dial(url)
	if err != nil {
		return err
	}
	c.client = client
	ecdsaKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}
	c.ethAddress = &c.authTransactor.From
	authTransactor := bind.NewKeyedTransactor(ecdsaKey)
	c.authTransactor = authTransactor

	oracleEmitter, err := oracleemitter.NewOracleEmitter(common.HexToAddress(oracleEmitterContractAddress), client)
	if err != nil {
		return err
	}
	aggregatorPlainSC, err := aggregator.NewAggregator(common.HexToAddress(aggregatorContractAddress), client)
	if err != nil {
		return err
	}
	c.oracleEmitter = &oracleemitter.OracleEmitterSession{
		Contract: oracleEmitter,
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
			GasLimit: 0,   // 0 automatically estimates gas limit
			GasPrice: nil, // nil automatically suggests gas price
			Context:  context.Background(),
		},
	}
	return nil
}

// // Balance returns the balance of the given ethereum address.
// func (c *EthereumClient) Balance(ctx context.Context, address string) (*big.Int, error) {
// 	ethereumAddress := common.HexToAddress(address)
// 	value, err := c.HttpClient.BalanceAt(ctx, ethereumAddress, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return value, nil
// }

// func (c *EthereumClient) SendTransaction(ctx context.Context, private_key, to string, amount int64) string {
// 	privateKey, err := crypto.HexToECDSA(private_key)
// 	if err != nil {
// 		c.Logger.Fatal("Failed to parse private key", err)
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		c.Logger.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey", err)
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	nonce, err := c.HttpClient.PendingNonceAt(ctx, fromAddress)
// 	if err != nil {
// 		c.Logger.Fatal("Failed to generate wallet nonce value", err)
// 	}

// 	value := big.NewInt(amount)
// 	gasLimit := uint64(21000) // in units
// 	gasPrice, err := c.HttpClient.SuggestGasPrice(ctx)
// 	if err != nil {
// 		c.Logger.Fatal("Failed to suggest new gas price", err)
// 	}

// 	toAddress := common.HexToAddress(to)
// 	var data []byte
// 	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

// 	chainID, err := c.HttpClient.NetworkID(ctx)
// 	if err != nil {
// 		c.Logger.Fatal("Failed to get network ID", err)
// 	}

// 	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
// 	if err != nil {
// 		c.Logger.Fatal("Failed to sign transaction", err)
// 	}

// 	err = c.HttpClient.SendTransaction(ctx, signedTx)
// 	if err != nil {
// 		c.Logger.Fatal("Failed to send signed transaction", err)
// 	}

// 	TxHash := signedTx.Hash().Hex()

// 	c.Logger.Info("Transaction sent: %s", TxHash)

// 	return TxHash
// }

func (c *EthereumClient) SubscribeOnOracleEvents(incomingEventsChan chan *oracleemitter.OracleEmitterNewOracleRequest) (event.Subscription, error) {
	requestsFilter := c.oracleEmitter.Contract.OracleEmitterFilterer
	subscription, err := requestsFilter.WatchNewOracleRequest(&bind.WatchOpts{
		Start:   nil, //last block
		Context: nil,
	}, incomingEventsChan)
	if err != nil {
		return nil, err
	}
	return subscription, err
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

	_, err := c.aggregator.CollectData(reqID, data, callbackAddress, callbackMethodID)
	if err != nil {
		return err
	}

	return nil
}

// Getting total stake in DioneStaking contract, this function could
// be used for storing the total stake and veryfing the stake tokens
// on new tasks
func (c *EthereumClient) GetTotalStake() (*big.Int, error) {
	totalStake, err := c.dioneStaking.TotalStake()
	if err != nil {
		return nil, err
	}
	return totalStake, nil
}

// Getting miner stake in DioneStaking contract, this function could
// be used for storing the miner's stake and veryfing the stake tokens
// on new tasks
func (c *EthereumClient) GetMinerStake(minerAddress common.Address) (*big.Int, error) {
	minerStake, err := c.dioneStaking.MinerStake(minerAddress)
	if err != nil {
		return nil, err
	}
	return minerStake, nil
}
