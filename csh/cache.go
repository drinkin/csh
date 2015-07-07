package csh

import (
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("Cache miss")

type Cache interface {
	Set(key string, value interface{}, exp ...time.Duration) error
	Get(key string) Value
	Del(key string) error
}

type Value interface {
	Int() (int, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Float64() (float64, error)
	String() (string, error)
	Bytes() ([]byte, error)
	Bool() (bool, error)
	JSON(obj interface{}) error
	Scan(obj interface{}) error
}

type Marshaler interface {
	MarshalCache() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalCache(data []byte) error
}

// Ensure all implement interface
var _ Value = (*RedisValue)(nil)
var _ Cache = (*Redis)(nil)
