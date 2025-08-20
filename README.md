# GBase - 企业级 Golang 微服务基础库

GBase 是一个功能完整的 Golang 微服务基础库，基于 Echo v4 框架构建，提供了构建现代 Web 服务所需的全套基础设施组件。

## ✨ 核心特性

- 🚀 **高性能 HTTP 服务器** - 基于 Echo v4，支持中间件和优雅关闭
- ⚙️ **统一配置管理** - 环境变量驱动的配置系统，支持多环境部署
- 📚 **自动 API 文档** - 内置 Swagger 文档生成和 UI 界面
- 🗄️ **数据库支持** - MySQL 连接池管理，支持主从分离
- 🔄 **Redis 集成** - Redis 客户端封装和连接管理
- 📝 **结构化日志** - 基于 slog 的高性能日志系统
- 💾 **智能缓存** - 内存和 Redis 双层缓存，支持并发安全
- 🛡️ **错误处理** - 统一的业务错误类型和 HTTP 错误处理
- 🔧 **类型绑定** - 自动请求参数绑定和验证
- 📦 **模块化设计** - 清晰的包结构，易于扩展和维护

## 📦 安装

```bash
go get github.com/just-hai/gbase
```

## 🚀 快速开始

### 基础 HTTP 服务

```go
package main

import (
    "github.com/just-hai/gbase/pkg/server"
    "github.com/just-hai/gbase/pkg/types"
)

func main() {
    // 创建服务器实例（自动加载环境变量配置）
    srv := server.New()

    // 创建 API 路由组
    api := srv.RouterGroup.Group("/api/v1", "API v1")
    
    // 注册路由
    api.GET("/users", getUsers, "获取用户列表")
    api.POST("/users", createUser, "创建用户")
    api.GET("/users/:id", getUserByID, "根据ID获取用户")

    // 启动服务器（支持优雅关闭）
    srv.Start()
}

// 处理器函数 - 支持自动参数绑定和类型转换
func getUsers(req *GetUsersRequest) (*GetUsersResponse, error) {
    // 业务逻辑
    users := []User{
        {ID: 1, Name: "张三", Email: "zhangsan@example.com"},
        {ID: 2, Name: "李四", Email: "lisi@example.com"},
    }
    
    return &GetUsersResponse{
        Users: users,
        Total: len(users),
    }, nil
}

func createUser(req *CreateUserRequest) (*User, error) {
    // 参数验证已自动完成
    user := &User{
        ID:    generateID(),
        Name:  req.Name,
        Email: req.Email,
    }
    
    // 保存到数据库...
    
    return user, nil
}

func getUserByID(req *GetUserByIDRequest) (*User, error) {
    if req.ID <= 0 {
        return nil, types.CodeBizError(400, "无效的用户ID")
    }
    
    // 查询数据库...
    user := &User{ID: req.ID, Name: "用户" + string(req.ID)}
    return user, nil
}

// 请求和响应结构体
type GetUsersRequest struct {
    Page     int    `query:"page" validate:"min=1" default:"1" desc:"页码"`
    PageSize int    `query:"page_size" validate:"min=1,max=100" default:"10" desc:"每页数量"`
    Keyword  string `query:"keyword" desc:"搜索关键词"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=50" desc:"用户名"`
    Email string `json:"email" validate:"required,email" desc:"邮箱地址"`
}

type GetUserByIDRequest struct {
    ID int `param:"id" validate:"required,min=1" desc:"用户ID"`
}

type User struct {
    ID    int    `json:"id" desc:"用户ID"`
    Name  string `json:"name" desc:"用户名"`
    Email string `json:"email" desc:"邮箱地址"`
}

type GetUsersResponse struct {
    Users []User `json:"users" desc:"用户列表"`
    Total int    `json:"total" desc:"总数量"`
}
```

## ⚙️ 配置管理

GBase 使用环境变量进行配置管理，支持以下配置项：

### 基础配置

```bash
# 应用基础信息
export APP_NAME="my-service"           # 应用名称
export APP_VERSION="1.0.0"            # 应用版本
export PORT=":8080"                    # 监听端口
export ENV="development"               # 环境：development/staging/prod
export SWAGGER="true"                  # 是否启用 Swagger 文档

# 日志配置
export LOG_LEVEL="INFO"                # 日志级别：DEBUG/INFO/WARN/ERROR
```

### 数据库配置

```bash
# MySQL 配置
export MYSQL_DSN="user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
export MYSQL_SLAVE_DSN="user:password@tcp(slave-host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
export MYSQL_MAX_OPEN_CONN="300"       # 最大连接数
export MYSQL_MAX_IDLE_CONN="100"       # 最大空闲连接数
export MYSQL_CONN_MAX_LIFETIME="10"    # 连接最大生存时间（分钟）
export MYSQL_TIMEOUT="5"               # 连接超时时间（秒）
export MYSQL_READ_TIMEOUT="30"         # 读取超时时间（秒）
export MYSQL_WRITE_TIMEOUT="30"        # 写入超时时间（秒）
```

### Redis 配置

```bash
# Redis 配置
export REDIS_ADDR="localhost:6379"     # Redis 地址
export REDIS_PW=""                     # Redis 密码
export REDIS_DB="0"                    # Redis 数据库编号
```

## 🗄️ 数据库使用

```go
import (
    "github.com/just-hai/gbase/pkg/db"
    "github.com/just-hai/gbase/pkg/config"
)

func initDB() {
    // 初始化数据库连接
    err := db.InitMysql(&config.Conf.Mysql)
    if err != nil {
        log.Fatal("数据库初始化失败:", err)
    }
    
    // 使用主库
    result := db.Master.Create(&user)
    
    // 使用从库（只读）
    var users []User
    db.Slave.Find(&users)
}
```

## 🔄 Redis 使用

```go
import (
    "github.com/just-hai/gbase/pkg/redis"
    "github.com/just-hai/gbase/pkg/config"
)

func initRedis() {
    // 初始化 Redis 连接
    redis.InitRedis(&config.Conf.Redis)
    
    // 使用 Redis
    redis.Client.Set(ctx, "key", "value", time.Hour)
    val := redis.Client.Get(ctx, "key").Val()
}
```

## 💾 缓存系统

GBase 提供了高性能的双层缓存系统：

```go
import "github.com/just-hai/gbase/pkg/cache"

// 内存缓存 - 适用于高频访问的小数据
func getUserFromMemCache(userID int) (*User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    return cache.OnceInMem(key, 5*time.Minute, func() (*User, error) {
        // 缓存未命中时的数据获取逻辑
        return getUserFromDB(userID)
    })
}

// Redis 缓存 - 适用于需要持久化的缓存数据
func getUserFromRedisCache(userID int) (*User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    return cache.OnceInRedis(key, time.Hour, func() (*User, error) {
        // 缓存未命中时的数据获取逻辑
        return getUserFromDB(userID)
    })
}
```

## 📝 日志系统

```go
import (
    "github.com/just-hai/gbase/pkg/log"
    "log/slog"
)

func initLogging() {
    // 初始化日志系统
    log.Init(log.LogOptions{
        Level: "INFO",
        Handle: func(time time.Time, msg, level string, attrs []slog.Attr) {
            // 自定义日志处理逻辑
            fmt.Printf("%s [%s] %s %+v\n", 
                time.Format("2006-01-02 15:04:05"), level, msg, attrs)
        },
    })
    
    // 使用结构化日志
    slog.With("user_id", 123, "action", "login").Info("用户登录成功")
    slog.Error("数据库连接失败", "error", err)
}
```

## 📚 API 文档

启用 Swagger 后，可通过以下地址访问：

- **Swagger JSON**: `http://localhost:8080/swagger.json`
- **Swagger UI**: `http://localhost:8080/swagger/`

API 文档会根据你的路由定义和结构体标签自动生成。

## 🏗️ 项目结构

```text
pkg/
├── cache/           # 缓存系统（内存缓存、Redis缓存）
│   ├── cache.go     # 分片内存缓存实现
│   └── once.go      # 单次执行缓存（防缓存击穿）
├── config/          # 配置管理
│   └── config.go    # 环境变量配置结构
├── db/              # 数据库连接管理
│   └── mysql.go     # MySQL 连接池和主从分离
├── log/             # 日志系统
│   ├── log.go       # 结构化日志实现
│   └── const.go     # 日志常量定义
├── redis/           # Redis 客户端
│   └── redis.go     # Redis 连接管理
├── server/          # HTTP 服务器
│   ├── handler/     # 请求处理器
│   │   ├── caller.go    # 处理器调用逻辑
│   │   ├── response.go  # 响应格式化
│   │   └── types.go     # 类型绑定和验证
│   ├── middleware/  # 中间件
│   │   └── request-logger.go # 请求日志中间件
│   ├── router/      # 路由管理
│   │   └── group.go     # 路由组和 Swagger 集成
│   └── server.go    # 服务器核心实现
├── swagger/         # API 文档生成
│   ├── generator.go # Swagger 文档生成器
│   ├── example.go   # 示例定义
│   └── ui/          # Swagger UI 静态文件
├── types/           # 公共类型定义
│   └── error.go     # 业务错误类型
└── utils/           # 工具函数
```

## 🔧 高级特性

### 自动参数绑定和验证

GBase 支持自动将 HTTP 请求参数绑定到结构体，并进行验证：

```go
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=1,max=200" desc:"文章标题"`
    Content string `json:"content" validate:"required,min=10" desc:"文章内容"`
    Tags    []string `json:"tags" validate:"max=10" desc:"标签列表"`
    UserID  int    `query:"user_id" validate:"required,min=1" desc:"用户ID"`
    Token   string `header:"Authorization" validate:"required" desc:"认证令牌"`
}
```

### 业务错误处理

```go
import "github.com/just-hai/gbase/pkg/types"

func someBusinessLogic() error {
    // 自定义错误码和消息
    return types.CodeBizError(1001, "用户不存在")
    
    // 通用业务错误
    return types.CustomBizError("操作失败，请稍后重试")
}
```

### 中间件支持

```go
// 自定义中间件
func AuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // 认证逻辑
            token := c.Request().Header.Get("Authorization")
            if token == "" {
                return types.CodeBizError(401, "未授权访问")
            }
            return next(c)
        }
    }
}

// 使用中间件
api := srv.RouterGroup.Group("/api/v1", "API v1")
api.Use(AuthMiddleware())
```

## 📋 最佳实践

1. **配置管理**: 使用环境变量管理不同环境的配置
2. **错误处理**: 使用统一的业务错误类型，提供清晰的错误信息
3. **日志记录**: 使用结构化日志，包含足够的上下文信息
4. **缓存策略**: 根据数据特性选择内存缓存或 Redis 缓存
5. **数据库连接**: 合理配置连接池参数，使用主从分离提高性能
6. **API 设计**: 遵循 RESTful 规范，提供完整的 API 文档

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进 GBase。

## 📄 许可证

MIT License