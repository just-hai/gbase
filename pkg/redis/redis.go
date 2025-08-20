package redis

import (
	"context"

	"log"

	"github.com/just-hai/gbase/pkg/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	// PoolSize 默认值 default: 10 * runtime.GOMAXPROCS(0)
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})
	if pingVal := rdb.Ping(context.TODO()); pingVal.Err() != nil {
		panic(pingVal.Err())
	} else {
		log.Printf("redis init success: %s", pingVal.Val())
	}
	Client = rdb
}
