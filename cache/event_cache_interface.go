package cache

import "github.com/Secured-Finance/dione/contracts/dioneOracle"

type EventCache interface {
	Store(key string, event interface{}) error
	GetOracleRequestEvent(key string) (*dioneOracle.DioneOracleNewOracleRequest, error)
	Delete(key string)
}
