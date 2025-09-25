// Package db -----------------------------
// @file      : mysql.go
// @author    : Carlos
// @contact   : 534994749@qq.com
// @time      : 2025/6/17 14:57
// -------------------------------------------
package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type Config struct {
	Dsn             string
	TablePrefix     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime time.Duration
	Debug           bool
}

// InitMysql 初始化mysql
func InitMysql(config *Config) (*gorm.DB, error) {
	gConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix, // 增加表前缀
			SingularTable: true,               // 使用单数表名
		},
	}
	// 显示所有SQL日志
	if config.Debug {
		gConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(config.Dsn), gConfig)
	if err != nil {
		return nil, err
	}
	// 获取连接池
	sqlDB, _ := db.DB()
	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	// 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	// 注：一定要比MYSQL的wait_timeout小，否则可能出现如下错误
	// wsarecv: An established connection was aborted by the software in your host machine.
	sqlDB.SetConnMaxLifetime(time.Second * config.ConnMaxLifeTime)
	return db, nil
}
