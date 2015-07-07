package csh_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/drinkin/csh/csh"
	"github.com/drinkin/di/rdis"
	"github.com/drinkin/shop/src/di/random"
	"github.com/stretchr/testify/assert"
)

type Example struct {
	UserId int64
}

type ExampleCache struct {
	UserId int64
}

func (ec *ExampleCache) MarshalCache() ([]byte, error) {
	return json.Marshal(ec)
}

func (ec *ExampleCache) UnmarshalCache(data []byte) error {
	return json.Unmarshal(data, ec)
}

var cacheTests = []struct {
	v   interface{}
	exp time.Duration
}{
	{v: 1},
	{v: false},
	{v: true, exp: time.Duration(1) * time.Hour},
	{v: 3, exp: time.Duration(1) * time.Second},
	{v: int64(-6)},
	{v: uint64(4)},
	{v: float64(4.341)},
	{v: "hello test"},
}

func resToType(expected interface{}, res csh.Value) (interface{}, error) {
	var o interface{}
	var err error
	switch expected.(type) {
	case int64:
		o, err = res.Int64()
	case int:
		o, err = res.Int()
	case uint64:
		o, err = res.Uint64()
	case float64:
		o, err = res.Float64()
	case bool:
		o, err = res.Bool()
	case string:
		o, err = res.String()
	default:
		panic(expected)
	}

	return o, err
}

func TestRedis(t *testing.T) {
	assert := assert.New(t)
	pool := rdis.NewPool("localhost:6379", "")
	c := csh.NewRedis(pool)

	for _, ct := range cacheTests {
		k := random.String(16)
		c.Set(k, ct.v, ct.exp)
		o, err := resToType(ct.v, c.Get(k))
		assert.NoError(err)
		assert.Equal(ct.v, o)

		{
			c.Del(k)
			_, err := resToType(ct.v, c.Get(k))
			assert.Equal(csh.ErrCacheMiss, err, "The value should be deleted")
		}

	}

	v := &Example{1}
	b, err := json.Marshal(v)
	assert.NoError(err)
	c.Set("d", b, time.Duration(1)*time.Second)

	var e Example
	err = c.Get("d").JSON(&e)
	assert.NoError(err)
	assert.EqualValues(1, e.UserId)

	ec := &ExampleCache{random.Int64(0, 1000)}
	key := random.String(10)
	c.Set(key, ec)
	var ec_val ExampleCache
	assert.NoError(c.Get(key).Scan(&ec_val))
	assert.EqualValues(ec, &ec_val)
}
