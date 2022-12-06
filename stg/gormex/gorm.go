package gormex

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Debug           bool
	SilentLog       bool
}

func InitGorm(dsn string) (db *gorm.DB, err error) {
	return InitGormWithConfig(DBConfig{
		DSN: dsn,
	})
}

func InitGormWithConfig(cfg DBConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if cfg.Debug {
		db = db.Debug()
	} else {
		level := logger.Warn
		if cfg.SilentLog {
			level = logger.Silent
		}
		db = db.Session(&gorm.Session{
			Logger: db.Logger.LogMode(level),
		})
	}

	return
}
