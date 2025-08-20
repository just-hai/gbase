package main

import (
	"fmt"
	"log"
	"time"

	"github.com/just-hai/gbase/pkg/cache"
)

// User 用户结构体
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Product 产品结构体
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	fmt.Println("=== 泛型缓存使用示例 ===\n")

	// 创建缓存实例
	c := cache.New()

	// 示例 1: 字符串类型缓存
	fmt.Println("1. 字符串类型缓存:")
	result1, err := cache.Get(c, "greeting", 5*time.Second, func() (string, error) {
		fmt.Println("  执行函数获取字符串...")
		time.Sleep(100 * time.Millisecond)
		return "Hello, World!", nil
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("  结果: %s (类型: %T)\n", result1, result1)
	}

	// 第二次调用，从缓存返回
	result1, _ = cache.Get(c, "greeting", 5*time.Second, func() (string, error) {
		fmt.Println("  这不应该被执行")
		return "Hello, World!", nil
	})
	fmt.Printf("  缓存结果: %s\n\n", result1)

	// 示例 2: 结构体类型缓存
	fmt.Println("2. 结构体类型缓存:")
	user, err := cache.Get(c, "user_123", 5*time.Second, func() (User, error) {
		fmt.Println("  执行函数获取用户信息...")
		time.Sleep(200 * time.Millisecond)
		return User{
			ID:    123,
			Name:  "张三",
			Email: "zhangsan@example.com",
		}, nil
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("  用户: %+v (类型: %T)\n\n", user, user)
	}

	// 示例 3: 切片类型缓存
	fmt.Println("3. 切片类型缓存:")
	products, err := cache.Get(c, "products", 5*time.Second, func() ([]Product, error) {
		fmt.Println("  执行函数获取产品列表...")
		time.Sleep(150 * time.Millisecond)
		return []Product{
			{ID: 1, Name: "笔记本电脑", Price: 5999.99},
			{ID: 2, Name: "无线鼠标", Price: 199.99},
			{ID: 3, Name: "机械键盘", Price: 899.99},
		}, nil
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("  产品列表: %+v (类型: %T)\n", products, products)
		fmt.Printf("  产品数量: %d\n\n", len(products))
	}

	// 示例 4: Map 类型缓存
	fmt.Println("4. Map 类型缓存:")
	config, err := cache.Get(c, "app_config", 5*time.Second, func() (map[string]interface{}, error) {
		fmt.Println("  执行函数获取配置...")
		time.Sleep(100 * time.Millisecond)
		return map[string]interface{}{
			"app_name":  "MyApp",
			"version":   "1.0.0",
			"debug":     true,
			"max_users": 1000,
			"timeout":   30.5,
		}, nil
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("  配置: %+v (类型: %T)\n\n", config, config)
	}

	// 示例 5: 整数类型缓存
	fmt.Println("5. 整数类型缓存:")
	count, err := cache.Get(c, "user_count", 5*time.Second, func() (int, error) {
		fmt.Println("  执行函数计算用户数量...")
		time.Sleep(80 * time.Millisecond)
		return 42, nil
	})
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("  用户数量: %d (类型: %T)\n\n", count, count)
	}

	// 示例 6: 类型安全演示
	fmt.Println("6. 类型安全演示:")

	// 这样使用是类型安全的，不需要类型断言
	userName := user.Name // 直接访问字段，无需类型断言
	userEmail := user.Email
	fmt.Printf("  用户名: %s, 邮箱: %s\n", userName, userEmail)

	// 遍历产品列表也是类型安全的
	fmt.Println("  产品详情:")
	for i, product := range products {
		fmt.Printf("    %d. %s - ¥%.2f\n", i+1, product.Name, product.Price)
	}

	fmt.Println("\n=== 示例结束 ===")
}
