package main

import (
	"fmt"
	"log"
	"time"

	"github.com/just-hai/gbase/pkg/cache"
)

// 模拟一个耗时的数据库查询
func fetchUserFromDB(userID string) (interface{}, error) {
	fmt.Printf("正在从数据库获取用户 %s 的信息...\n", userID)

	// 模拟数据库查询延迟
	time.Sleep(500 * time.Millisecond)

	// 返回模拟的用户数据
	return map[string]interface{}{
		"id":    userID,
		"name":  fmt.Sprintf("User_%s", userID),
		"email": fmt.Sprintf("user_%s@example.com", userID),
	}, nil
}

// 模拟一个可能失败的 API 调用
func fetchDataFromAPI(key string) (interface{}, error) {
	fmt.Printf("正在从 API 获取数据 %s...\n", key)

	// 模拟 API 调用延迟
	time.Sleep(300 * time.Millisecond)

	// 模拟偶尔的失败
	if key == "error_key" {
		return nil, fmt.Errorf("API 调用失败")
	}

	return fmt.Sprintf("API_Data_%s", key), nil
}

func main() {
	// 创建缓存实例
	userCache := cache.New()

	fmt.Println("=== 缓存使用示例 ===\n")

	// 示例 1: 基本缓存使用
	fmt.Println("1. 基本缓存使用:")

	// 第一次调用 - 会执行函数
	start := time.Now()
	result, err := userCache.Get("user_123", 5*time.Second, func() (interface{}, error) {
		return fetchUserFromDB("123")
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("第一次调用结果: %+v (耗时: %v)\n", result, time.Since(start))
	}

	// 第二次调用 - 从缓存返回
	start = time.Now()
	result, err = userCache.Get("user_123", 5*time.Second, func() (interface{}, error) {
		return fetchUserFromDB("123")
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("第二次调用结果: %+v (耗时: %v)\n", result, time.Since(start))
	}

	fmt.Println()

	// 示例 2: 并发调用相同 key
	fmt.Println("2. 并发调用相同 key (singleflight 效果):")

	start = time.Now()
	done := make(chan bool, 3)

	// 启动 3 个并发请求
	for i := 0; i < 3; i++ {
		go func(id int) {
			result, err := userCache.Get("concurrent_user", 5*time.Second, func() (interface{}, error) {
				return fetchUserFromDB("concurrent")
			})
			if err != nil {
				log.Printf("协程 %d 错误: %v", id, err)
			} else {
				fmt.Printf("协程 %d 结果: %+v\n", id, result)
			}
			done <- true
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < 3; i++ {
		<-done
	}
	fmt.Printf("并发调用总耗时: %v\n", time.Since(start))

	fmt.Println()

	// 示例 3: 错误处理
	fmt.Println("3. 错误处理:")

	result, err = userCache.Get("error_data", 5*time.Second, func() (interface{}, error) {
		return fetchDataFromAPI("error_key")
	})
	if err != nil {
		fmt.Printf("预期的错误: %v\n", err)
	}

	// 错误不会被缓存，再次调用会重新执行
	result, err = userCache.Get("error_data", 5*time.Second, func() (interface{}, error) {
		return fetchDataFromAPI("success_key")
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("成功调用结果: %v\n", result)
	}

	fmt.Println()

	// 示例 4: 缓存过期
	fmt.Println("4. 缓存过期:")

	// 设置很短的过期时间
	result, _ = userCache.Get("short_lived", 1*time.Second, func() (interface{}, error) {
		return "第一次数据", nil
	})
	fmt.Printf("第一次结果: %v\n", result)

	// 等待过期
	fmt.Println("等待缓存过期...")
	time.Sleep(1500 * time.Millisecond)

	result, _ = userCache.Get("short_lived", 1*time.Second, func() (interface{}, error) {
		return "第二次数据", nil
	})
	fmt.Printf("过期后结果: %v\n", result)

	fmt.Println()

	// 示例 5: 缓存管理
	fmt.Println("5. 缓存管理:")

	fmt.Printf("当前缓存大小: %d\n", userCache.Size())

	// 清理过期项
	userCache.CleanExpired()
	fmt.Printf("清理过期项后大小: %d\n", userCache.Size())

	// 删除特定项
	userCache.Delete("user_123")
	fmt.Printf("删除特定项后大小: %d\n", userCache.Size())

	// 清空所有缓存
	userCache.Clear()
	fmt.Printf("清空后大小: %d\n", userCache.Size())

	fmt.Println("\n=== 示例结束 ===")
}
