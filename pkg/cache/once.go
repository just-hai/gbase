package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/just-hai/gbase/pkg/redis"
	"golang.org/x/sync/singleflight"
)

var (
	onceSg      singleflight.Group
	onceCache   = NewShardedCache(16, 5*time.Second)
	ErrOnceType = errors.New("cache item is not of the correct type")
)

func checkMemExist[T any](key string) (T, error) {
	var zero T
	val, ok := onceCache.Get(key)
	if !ok {
		return zero, errors.New("item not found")
	}
	if value, ok := val.(T); ok {
		return value, nil
	} else {
		onceCache.Delete(key)
		return zero, ErrOnceType
	}
}

func OnceInMem[T any](key string, duration time.Duration, fn func() (T, error)) (T, error) {
	var zero T

	if key == "" {
		return zero, errors.New("key cannot be empty")
	}

	if item, err := checkMemExist[T](key); err == nil {
		return item, nil
	}

	result, err, _ := onceSg.Do(key, func() (interface{}, error) {
		if item, err := checkMemExist[T](key); err == nil {
			return item, nil
		}
		value, err := fn()
		if err != nil {
			return zero, err
		}
		onceCache.Set(key, value, duration)
		return value, nil
	})

	if err != nil {
		return zero, err
	}
	if value, ok := result.(T); ok {
		return value, nil
	}
	return zero, ErrOnceType
}

func checkRedisExist[T any](key string) (T, error) {
	zero := new(T)
	result := redis.Client.Get(context.TODO(), key)
	if result.Err() != nil {
		return *zero, result.Err()
	}

	val := result.Val()
	if val == "" {
		return *zero, errors.New("item not found")
	}
	err := json.Unmarshal([]byte(val), zero)
	return *zero, err
}

func OnceInRedis[T any](key string, duration time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	if redis.Client == nil {
		return zero, errors.New("redis client not initialized")
	}
	if key == "" {
		return zero, errors.New("key cannot be empty")
	}

	if item, err := checkRedisExist[T](key); err == nil {
		return item, nil
	}

	result, err, _ := onceSg.Do(key, func() (interface{}, error) {
		if item, err := checkRedisExist[T](key); err == nil {
			return item, nil
		}

		value, err := fn()
		if err != nil {
			return zero, err
		}
		by, err := json.Marshal(value)
		if err != nil {
			return zero, err
		}
		if err := redis.Client.SetEx(context.TODO(), key, string(by), duration).Err(); err != nil {
			return zero, err
		}
		return value, nil
	})

	if err != nil {
		return zero, err
	}

	if value, ok := result.(T); ok {
		return value, nil
	}

	return zero, ErrOnceType
}
