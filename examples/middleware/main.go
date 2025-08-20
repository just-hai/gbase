package main

import (
	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/server"
	"github.com/labstack/echo/v4"
)

func main() {
	// 创建服务器实例
	srv := server.NewWithConfig(&config.Config{
		Base: config.Base{
			AppName: "middleware-example",
			Port:    ":8081",
			Swagger: true,
		},
	})

	// 创建需要认证的路由组
	authAPI := srv.RouterGroup.Group("/api/auth", "authenticated")

	// 注册需要认证的路由
	authAPI.GET("/profile", authMiddleware, getProfile, "获取用户资料")
	authAPI.PUT("/profile", authMiddleware, updateProfile, "更新用户资料")

	// 启动服务器
	srv.Start()
}

// 认证中间件
func authMiddleware(ctx echo.Context) *User {
	// 这里可以从 Header 中获取 token 并验证
	// token := ctx.Request().Header.Get("Authorization")

	// 模拟认证成功，返回用户信息
	return &User{
		ID:   1,
		Name: "认证用户",
	}
}

// 获取用户资料
func getProfile(user *User) (*UserProfile, error) {
	return &UserProfile{
		ID:    user.ID,
		Name:  user.Name,
		Email: "user@example.com",
		Phone: "13800138000",
	}, nil
}

// 更新用户资料
func updateProfile(user *User, req *UpdateProfileRequest) (*UserProfile, error) {
	return &UserProfile{
		ID:    user.ID,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}, nil
}

// 数据结构
type User struct {
	ID   int    `json:"id" desc:"用户ID"`
	Name string `json:"name" desc:"用户名"`
}

type UserProfile struct {
	ID    int    `json:"id" desc:"用户ID"`
	Name  string `json:"name" desc:"用户名"`
	Email string `json:"email" desc:"邮箱"`
	Phone string `json:"phone" desc:"手机号"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" desc:"用户名"`
	Email string `json:"email" desc:"邮箱"`
	Phone string `json:"phone" desc:"手机号"`
}
