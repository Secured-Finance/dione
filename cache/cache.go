package cache

import (
	"errors"
	"time"
)

var ErrNilValue = errors.New("value is empty")

type Cache interface {
	Store(key string, value interface{}) error
	StoreWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string, value interface{}) error
	Delete(key string)
}
