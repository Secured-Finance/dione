package store

import "math/big"

// TODO: specify store for staking mechanism
type StakeTokenInfo struct {
	TotalTokensStaked *big.Int
	NodeTokensStaked  *big.Int
}
