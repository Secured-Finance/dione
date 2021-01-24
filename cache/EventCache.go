package cache

type EventCache interface {
	Store(key string, event interface{}) error
	GetOracleRequestEvent(key string) (*oracleEmitter.OracleEmitterNewOracleRequest, error)
	Delete(key string)
}