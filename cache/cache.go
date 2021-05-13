package cache

import "errors"

var ErrNilValue = errors.New("value is empty")

type Cache interface {
	Store(key string, value interface{}) error
	Get(key string, value interface{}) error
	Delete(key string)
}
