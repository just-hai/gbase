package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/redis"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Time int64
}

func getUser() user {
	return user{
		ID:   1,
		Name: "test",
		Time: time.Now().Unix(),
	}
}

func getUsers() []user {
	return []user{
		{ID: 1, Name: "张三", Time: time.Now().Unix()},
		{ID: 2, Name: "李四", Time: time.Now().Unix()},
	}
}

func getUserMap() map[int]*user {
	return map[int]*user{
		1: {ID: 1, Name: "张三", Time: time.Now().Unix()},
		2: {ID: 2, Name: "李四", Time: time.Now().Unix()},
	}
}

func TestOnceInMem(t *testing.T) {
	for i := range 10 {
		user, err := OnceInMem("test", 5*time.Second, func() (user, error) {
			fmt.Println("  执行函数获取用户信息...")
			return getUser(), nil
		})
		if i > 4 {
			time.Sleep(4 * time.Second)
		}
		fmt.Println(user, err)
	}
}

func TestOnceInMemInt(t *testing.T) {
	for i := range 10 {
		user, err := OnceInMem("test", 5*time.Second, func() (int64, error) {
			fmt.Println("  执行函数获取用户信息...")
			return time.Now().Unix(), nil
		})
		if i > 4 {
			time.Sleep(4 * time.Second)
		}
		fmt.Println(user, err)
	}
}

func TestOnceInRedis(t *testing.T) {
	InitRedis()
	for i := range 10 {
		user, err := OnceInRedis("test", 5*time.Second, func() ([]user, error) {
			fmt.Println("  执行函数获取用户信息...")
			return getUsers(), nil
		})
		if i > 4 {
			time.Sleep(4 * time.Second)
		}
		fmt.Println(user, err)
	}
}

func TestOnceInRedisMap(t *testing.T) {
	InitRedis()
	for i := range 10 {
		user, err := OnceInRedis("test", 5*time.Second, func() (map[int]*user, error) {
			fmt.Println("  执行函数获取用户信息...")
			return getUserMap(), nil
		})
		if i > 4 {
			time.Sleep(4 * time.Second)
		}
		fmt.Println(user, err)
	}
}

func InitRedis() {
	config.Conf.Redis.Addr = "122.112.252.179:6379"
	config.Conf.Redis.Password = "20RG8Yh72E"
	redis.Init()
}
