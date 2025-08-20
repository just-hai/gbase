package main

import (
	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/server"
	"github.com/just-hai/gbase/pkg/types"
	"github.com/labstack/echo/v4"
)

func main() {
	// 创建服务器实例
	srv := server.NewWithConfig(&config.Config{
		Base: config.Base{
			AppName: "basic-example",
			Port:    ":8080",
			Swagger: true,
		},
	})

	// 创建 API 路由组
	api := srv.RouterGroup.Group("/api", "api")

	// 注册路由
	api.GET("/users", getUsers, "获取用户列表")
	api.GET("/users/:id", getUserByID, "根据ID获取用户")
	api.POST("/users", createUser, "创建用户")

	// 启动服务器
	srv.Start()
}

// 获取用户列表
func getUsers() ([]User, error) {
	return []User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
	}, nil
}

// 根据ID获取用户
func getUserByID(ctx echo.Context) (*User, error) {
	// 这里可以从路径参数获取ID
	// id := ctx.Param("id")
	return &User{
		ID:    1,
		Name:  "张三",
		Email: "zhangsan@example.com",
	}, nil
}

// 创建用户
func createUser(req *CreateUserRequest) (*User, error) {
	if req.Name == "" {
		return nil, types.CustomBizError("用户名不能为空")
	}

	return &User{
		ID:    3,
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

// 数据结构定义
type User struct {
	ID    int    `json:"id" desc:"用户ID"`
	Name  string `json:"name" desc:"用户名"`
	Email string `json:"email" desc:"邮箱地址"`
}

type CreateUserRequest struct {
	Name  string `json:"name" desc:"用户名"`
	Email string `json:"email" desc:"邮箱地址"`
}
