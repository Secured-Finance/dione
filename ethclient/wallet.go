package ethclient

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

// // Balance returns the balance of the given ethereum address.
func (c *EthereumClient) Balance(ctx context.Context, address string) (*big.Int, error) {
	ethereumAddress := common.HexToAddress(address)
	value, err := c.client.BalanceAt(ctx, ethereumAddress, nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *EthereumClient) SendTransaction(ctx context.Context, private_key, to string, amount int64) string {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		logrus.Fatal("Failed to parse private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		logrus.Fatal("Failed to generate wallet nonce value", err)
	}

	value := big.NewInt(amount)
	gasLimit := uint64(21000) // in units
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		logrus.Fatal("Failed to suggest new gas price", err)
	}

	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := c.client.NetworkID(ctx)
	if err != nil {
		logrus.Fatal("Failed to get network ID", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		logrus.Fatal("Failed to sign transaction", err)
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		logrus.Fatal("Failed to send signed transaction", err)
	}

	TxHash := signedTx.Hash().Hex()

	logrus.Info("Transaction sent: %s", TxHash)

	return TxHash
}
