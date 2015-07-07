package csh

import (
	"encoding"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	pool *redis.Pool
}

func NewRedis(pool *redis.Pool) *Redis {
	return &Redis{
		pool: pool,
	}
}

func (r *Redis) Set(key string, value interface{}, exp ...time.Duration) error {
	conn := r.pool.Get()
	defer conn.Close()

	switch v := value.(type) {
	case Marshaler:
		data, err := v.MarshalCache()
		if err != nil {
			return err
		}

		value = data
	case encoding.BinaryMarshaler:
		data, err := v.MarshalBinary()
		if err != nil {
			return err
		}

		value = data
	}

	args := redis.Args{}.Add(key).Add(value)
	if hasExpire(exp) {
		ms := int(exp[0] / time.Millisecond)
		args = args.Add("PX").Add(ms)
	}

	_, err := conn.Do("SET", args...)
	return err
}

func (r *Redis) Get(key string) Value {
	conn := r.pool.Get()
	defer conn.Close()

	res, err := conn.Do("GET", key)
	if res == nil {
		return &RedisValue{nil, ErrCacheMiss}
	}
	return &RedisValue{res, err}
}

func (r *Redis) Del(key string) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
