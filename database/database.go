package database

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2023/11/10 23:48
 * @file: database.go
 * @description: 数据库连接
 */

var db *gorm.DB

// Database DataBaseConfig GORM
type Database struct {
	DSN         string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime int
	MaxIdleTime int
}

func (d *Database) New() *gorm.DB {

	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // Log level
			Colorful:                  false,       // 禁用彩色打印
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			ParameterizedQueries:      true,        // 启用参数化查询
		},
	)

	db, err := gorm.Open(mysql.Open(d.DSN), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})

	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
		return nil
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping database: %s", err))
		return nil
	}

	var result int
	err = sqlDB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		panic(fmt.Errorf("failed to query database: %s", err))
		return nil
	}

	// 设置连接池大小
	sqlDB.SetMaxOpenConns(d.MaxOpenConn)
	sqlDB.SetMaxIdleConns(d.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(d.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(d.MaxIdleTime) * time.Second)

	return db
}

func GetConnection() (db *gorm.DB) {
	return db
}
