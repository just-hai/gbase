package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/log"
	"github.com/just-hai/gbase/pkg/server"
	"github.com/just-hai/gbase/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// "net/http"

// "github.com/hai/go-base/swagger"

func main() {
	//server := server.NewWithConfig(&config.Config{Base: config.Base{AppName: "swagger-demo", Port: ":8081"}})

	log.Init(log.LogOptions{
		Level: "INFO",
		Handle: func(time time.Time, msg, level string, source *slog.Source, attrs []slog.Attr) {
			fmt.Printf("%s %s %s %+v[%d] %+v\n", time.Format("2006-01-02 15:04:05"), level, msg, source.Function, source.Line, attrs)
		},
	})
	slog.Info("start server")

	config.Conf.AppName = "swagger-demo"
	server := server.New()
	server.Echo.Use(middleware.CORS())
	g := server.RouterGroup.Group("/api", "api")
	g.GET("/users", token, getUsers, "用户列表")
	g.GET("/users2", getUsers2, "用户列表2")
	server.Start()
	
}

func getUsers(req *User) error {
	//panic("panic hai")
	return types.CodeBizError(123, "错误")
}

func getUsers2(req *User) (User, error) {

	return User{
		ID:   1,
		Name: "张三",
	}, nil
	// return []User{
	// 	{ID: 1, Name: "张三"},
	// 	{ID: 2, Name: "李四"},
	// }, nil
}

func token(ctx echo.Context) *User {
	return &User{
		ID:   111,
		Name: "张三11",
	}
}

type User struct {
	ID    int    `json:"id" desc:"用户ID" validate:"required"`
	Name  string `json:"name" desc:"用户名称" validate:"max=10"`
	Age   int    `validate:"required,max=10" `
	Users []User2
}

type User2 struct {
	ID   int    `json:"id" desc:"用户ID2"`
	Name string `json:"name" desc:"用户名称2"`
}
