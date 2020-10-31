package types

import (
	"errors"
	"math/big"

	big2 "github.com/filecoin-project/go-state-types/big"
	validation "github.com/go-ozzo/ozzo-validation"
)

var EmptyInt = BigInt{}

type BigInt = big2.Int

func NewInt(i uint64) BigInt {
	return BigInt{Int: big.NewInt(0).SetUint64(i)}
}

func BigFromBytes(b []byte) BigInt {
	i := big.NewInt(0).SetBytes(b)
	return BigInt{Int: i}
}

func ValidateBigInt(i *big.Int) validation.RuleFunc {
	return func(value interface{}) error {
		bigInt := i.IsInt64()
		if !bigInt {
			return errors.New("expected big integer")
		}
		return nil
	}
}
