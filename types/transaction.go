package types

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	Hash      []byte
	Timestamp int64
	Data      []byte
}

func CreateTransaction(data []byte) *Transaction {
	timestamp := time.Now().Unix()
	encodedData := hex.EncodeToString(data)
	hash := crypto.Keccak256([]byte(fmt.Sprintf("%d_%s", timestamp, encodedData)))
	return &Transaction{
		Hash:      hash,
		Timestamp: timestamp,
		Data:      data,
	}
}
