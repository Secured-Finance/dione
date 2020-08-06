package rpcclient

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"

	"github.com/Secured-Finance/p2p-oracle-node/contracts/aggregator"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-log"
)

type EthereumClient struct {
	HttpClient *ethclient.Client
	WsClient   *ethclient.Client
	Logger     *log.ZapEventLogger
}

type Ethereum interface {
	Connect(context.Context, string) error
	Balance(context.Context, string) (*big.Int, error)
	SubscribeOnSmartContractEvents(context.Context, string)
	GenerateAddressFromPrivateKey(string) string
	SendTransaction(string, string, int64) string
	createKeyStore(string) string
	importKeyStore(string, string) string
}

func NewEthereumClient() *EthereumClient {
	ethereumClient := &EthereumClient{
		Logger: log.Logger("rendezvous"),
	}

	return ethereumClient
}

func (c *EthereumClient) Connect(ctx context.Context, url string, connectionType string) error {
	client, err := ethclient.Dial(url)
	if err != nil {
		c.Logger.Fatal(err)
	}
	if connectionType == "websocket" {
		c.WsClient = client
	} else {
		c.HttpClient = client
	}
	return nil
}

// Balance returns the balance of the given ethereum address.
func (c *EthereumClient) Balance(ctx context.Context, address string) (*big.Int, error) {
	ethereumAddress := common.HexToAddress(address)
	value, err := c.HttpClient.BalanceAt(ctx, ethereumAddress, nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *EthereumClient) SendTransaction(private_key, to string, amount int64) string {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		c.Logger.Fatal("Failed to parse private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		c.Logger.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.HttpClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		c.Logger.Fatal("Failed to generate wallet nonce value", err)
	}

	value := big.NewInt(amount)
	gasLimit := uint64(21000) // in units
	gasPrice, err := c.HttpClient.SuggestGasPrice(context.Background())
	if err != nil {
		c.Logger.Fatal("Failed to suggest new gas price", err)
	}

	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := c.HttpClient.NetworkID(context.Background())
	if err != nil {
		c.Logger.Fatal("Failed to get network ID", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		c.Logger.Fatal("Failed to sign transaction", err)
	}

	err = c.HttpClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		c.Logger.Fatal("Failed to send signed transaction", err)
	}

	TxHash := signedTx.Hash().Hex()

	c.Logger.Info("Transaction sent: %s", TxHash)

	return TxHash
}

func (c *EthereumClient) GenerateAddressFromPrivateKey(private_key string) string {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		c.Logger.Fatal("Failed to generate private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		c.Logger.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	c.Logger.Info(hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address
}

func (c *EthereumClient) SubscribeOnSmartContractEvents(ctx context.Context, address string) {
	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := c.WsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		c.Logger.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			c.Logger.Fatal(err)
		case vLog := <-logs:
			c.Logger.Info(vLog) // pointer to event log
		}
	}

}

func (c *EthereumClient) createKeyStore(password string) string {
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		c.Logger.Fatal("Failed to create new keystore", err)
	}

	return account.Address.Hex()
}

func (c *EthereumClient) importKeyStore(filePath string, password string) string {
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.Logger.Fatal("Failed to read keystore file", err)
	}

	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		c.Logger.Fatal("Failed to import keystore", err)
	}

	return account.Address.Hex()
}

func (c *EthereumClient) SendConsensusValues(private_key string, smartContractAddress string, reqID *big.Int, data string, callbackAddress common.Address, callbackMethodID [4]byte) string {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		c.Logger.Fatal("Failed to generate private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		c.Logger.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	c.Logger.Info(hexutil.Encode(publicKeyBytes)[4:])

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.HttpClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		c.Logger.Fatal(err)
	}

	gasPrice, err := c.HttpClient.SuggestGasPrice(context.Background())
	if err != nil {
		c.Logger.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress(smartContractAddress)
	contract, err := aggregator.NewSmartcontracts(address, c.HttpClient)
	if err != nil {
		c.Logger.Fatal(err)
	}

	tx, err := contract.CollectData(auth, reqID, data, callbackAddress, callbackMethodID)
	if err != nil {
		c.Logger.Fatal(err)
	}

	txHash := tx.Hash().Hex()

	return txHash
}
