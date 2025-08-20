package db

import (
	"testing"

	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/log"
)

func TestMysql(t *testing.T) {
	config.Conf.Mysql = config.Mysql{
		DSN:             "cailianpress_dba:Cailianpress_888@tcp(119.3.42.116:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
		SlaveDSN:        "cailianpress_dba:Cailianpress_888@tcp(119.3.42.116:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
		ConnMaxLifetime: 10,
		MaxIdleConn:     100,
		MaxOpenConn:     300,
		Timeout:         5,
		WriteTimeout:    30,
		ReadTimeout:     30,
	}

	log.Init(log.LogOptions{Level: "DEBUG"})

	Init()

	var hais []hai
	Client.Table("hai_test").Limit(10).Find(&hais)
	Client.Table("hai_test").Create(&hai{Name: "zhanghai", Age: 18, Address: "beijing", Job: "engineer", Descr: "good", Column7: "column7", Column8: "column8"})
}

type hai struct {
	Id      int    `gorm:"column:id"`
	Name    string `gorm:"column:name"`
	Age     int    `gorm:"column:age"`
	Address string `gorm:"column:address"`
	Job     string `gorm:"column:job"`
	Descr   string `gorm:"column:descr"`
	Column7 string `gorm:"column:column_7"`
	Column8 string `gorm:"column:column_8"`
}
