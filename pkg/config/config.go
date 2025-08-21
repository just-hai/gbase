package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type (
	Base struct {
		AppName    string `envconfig:"APP_NAME" default:"just-hai"` // 服务名
		AppVersion string `envconfig:"APP_VERSION" default:"1.0.0"`
		Port       string `envconfig:"PORT"`    // 端口
		Env        string `envconfig:"ENV"`     // 环境信息
		Swagger    bool   `envconfig:"SWAGGER"` // 是否开启swagger
		CorHosts   string // 跨域域名
	}

	Log struct {
		Level  string `default:"INFO"` // 日志级别
	}

	Redis struct {
		Addr     string `envconfig:"REDIS_ADDR"`
		Password string `envconfig:"REDIS_PW"`
		DB       int    `envconfig:"REDIS_DB"`
	}

	Mysql struct {
		DSN             string `envconfig:"MYSQL_DSN"`
		SlaveDSN        string `envconfig:"MYSQL_SLAVE_DSN"`
		MaxOpenConn     int    `envconfig:"MYSQL_MAX_OPEN_CONN" default:"300"`
		MaxIdleConn     int    `envconfig:"MYSQL_MAX_IDLE_CONN" default:"100"`
		ConnMaxLifetime int    `envconfig:"MYSQL_CONN_MAX_LIFETIME" default:"10"`
		Timeout         int    `envconfig:"MYSQL_TIMEOUT" default:"5"`
		ReadTimeout     int    `envconfig:"MYSQL_READ_TIMEOUT" default:"30"`
		WriteTimeout    int    `envconfig:"MYSQL_WRITE_TIMEOUT" default:"30"`
	}

	Config struct {
		Base
		Log
		Redis
		Mysql
	}
)

var Conf Config

func IsProd() bool {
	return Conf.Base.Env == "prod"
}

func init() {
	err := envconfig.Process("", &Conf)
	if err != nil {
		log.Fatal("init conf error:", err)
	}
}
