package db

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	gmysql "github.com/go-sql-driver/mysql"
	"github.com/just-hai/gbase/pkg/config"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var Client *DBCache

func Init() {
	var err error
	gormLogger := GormLogger{logger.Default.LogMode(logger.Info)}
	gdb, err := gorm.Open(mysql.Open(timeOutDSN(config.Conf.Mysql.DSN)), &gorm.Config{
		PrepareStmt: true,
		Logger:      &gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		QueryFields: true,
	})

	if err != nil {
		e := fmt.Sprintf("gdb conection set up failed, %s\n", err.Error())
		slog.Error(e)
		panic(e)
	}

	if config.Conf.Mysql.SlaveDSN != "" {
		err := gdb.Use(dbresolver.Register(dbresolver.Config{
			Replicas:          []gorm.Dialector{mysql.Open(timeOutDSN(config.Conf.Mysql.SlaveDSN))},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}))
		if err != nil {
			e := fmt.Sprintf("gdb slave conection set up failed, %s\n", err.Error())
			slog.Error(e)
			panic(e)
		}
	}

	Client = &DBCache{gdb}
	db, err := gdb.DB()
	if err != nil {
		e := fmt.Sprintf("get db failed, %s\n", err.Error())
		slog.Error(e)
		panic(e)
	}

	db.SetMaxOpenConns(config.Conf.Mysql.MaxOpenConn)
	db.SetMaxIdleConns(config.Conf.Mysql.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(config.Conf.Mysql.ConnMaxLifetime) * time.Second)

	if db, err := gdb.Clauses(dbresolver.Read).DB(); err != nil {
		log.Fatalf("db slave ping err, %s\n", err.Error())
	} else {
		if err := db.Ping(); err != nil {
			log.Fatalf("db slave ping err, %s\n", err.Error())
		}
	}

	if db, err := gdb.Clauses(dbresolver.Write).DB(); err != nil {
		log.Fatalf("db slave ping err, %s\n", err.Error())
	} else {
		if err := db.Ping(); err != nil {
			log.Fatalf("db slave ping err, %s\n", err.Error())
		}
	}

	slog.Info("db init success")
}

func timeOutDSN(dsn string) string {
	cfg, err := gmysql.ParseDSN(dsn)
	if err != nil {
		return ""
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = time.Duration(config.Conf.Mysql.Timeout) * time.Second
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = time.Duration(config.Conf.Mysql.ReadTimeout) * time.Second
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = time.Duration(config.Conf.Mysql.WriteTimeout) * time.Second
	}
	return cfg.FormatDSN()
}

type GormLogger struct {
	logger.Interface
}

func (l *GormLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	var (
		forceLog bool
		log      *slog.Logger
	)

	if _func := ctx.Value("func"); _func != nil {
		_func, _ := _func.(string)
		if _func != "" {
			log = slog.With("func", _func)
		}
		forceLog = true
	}
	if log == nil {
		log = slog.Default()
	}

	sql, rowsAffected := fc()

	sql = strings.TrimLeft(sql, " ")
	sqlPrefix := sql
	if len(sqlPrefix) > 20 {
		sqlPrefix = strings.ToUpper(sql[:20])
	} else {
		sqlPrefix = strings.ToUpper(sql)
	}

	elapsed := time.Since(begin)
	str := fmt.Sprintf("SQL: %s | RowsAffected: %d | Time: %v | Error: %v", sql, rowsAffected, elapsed, err)
	if err != nil {
		log.Error(str)
	} else if forceLog {
		log.Info(str)
	} else if strings.HasPrefix(sqlPrefix, "INSERT") || strings.HasPrefix(sqlPrefix, "UPDATE") ||
		strings.HasPrefix(sqlPrefix, "DELETE") {
		log.Info(str)
	} else if !config.IsProd() {
		log.Debug(str)
	}
}

type DBCache struct {
	*gorm.DB
}

func (dbCache *DBCache) Tag(_func string) *gorm.DB {
	return dbCache.DB.WithContext(context.WithValue(context.TODO(), "func", _func))
}
