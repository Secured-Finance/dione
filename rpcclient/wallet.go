package rpcclient

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

func GenerateEthWalletAddressFromPrivateKey(private_key string) common.Address {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		logrus.Fatal("Failed to generate private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	logrus.Info(hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address
}

// Convert common.Address type into string for ethereum wallet
func EthWalletToString(ethWallet common.Address) string {
	return ethWallet.Hex()
}
