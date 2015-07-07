package csh

import (
	"encoding"
	"encoding/json"
	"errors"

	"github.com/garyburd/redigo/redis"
)

type RedisValue struct {
	Raw interface{}
	Err error
}

func (v *RedisValue) Int() (int, error) {
	return redis.Int(v.Raw, v.Err)
}

func (v *RedisValue) Int64() (int64, error) {
	return redis.Int64(v.Raw, v.Err)
}

func (v *RedisValue) Uint64() (uint64, error) {
	return redis.Uint64(v.Raw, v.Err)
}

func (v *RedisValue) Float64() (float64, error) {
	return redis.Float64(v.Raw, v.Err)
}

func (v *RedisValue) String() (string, error) {

	return redis.String(v.Raw, v.Err)
}

func (v *RedisValue) Bytes() ([]byte, error) {
	return redis.Bytes(v.Raw, v.Err)
}

func (v *RedisValue) Bool() (bool, error) {
	return redis.Bool(v.Raw, v.Err)
}

func (v *RedisValue) JSON(obj interface{}) error {
	b, err := v.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func (v *RedisValue) Scan(obj interface{}) error {
	b, err := v.Bytes()
	if err != nil {
		return err
	}

	switch v := obj.(type) {
	case Unmarshaler:
		return v.UnmarshalCache(b)
	case encoding.BinaryUnmarshaler:
		return v.UnmarshalBinary(b)

	}
	return errors.New("Obj does not implement UnmarshalCache or UnmarshalBinary")
}
