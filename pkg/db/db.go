// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package db

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLOptions 定义 MySQL 数据库的选项.
type MySQLOptions struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
}

// DSN 从 MySQLOptions 返回 DSN.
func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")
}

// NewMySQL 使用给定的选项创建一个新的 gorm 数据库实例.
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	logLevel := logger.Silent
	if opts.LogLevel != 0 {
		logLevel = logger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns 设置到数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime 设置连接可重用的最长时间
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns 设置空闲连接池的最大连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

// PostgresOptions 定义 PostgreSQL 数据库的选项.
type PostgresOptions struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	SSLMode               string
	Timezone              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
}

// DSN 从 PostgresOptions 返回 DSN.
func (o *PostgresOptions) DSN() string {
	host := o.Host
	port := o.Port
	if port <= 0 {
		port = 5432
	}
	if h, p, err := net.SplitHostPort(o.Host); err == nil {
		host = h
		if pn, err := strconv.Atoi(p); err == nil {
			port = pn
		}
	}
	sslMode := o.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	tz := o.Timezone
	if tz == "" {
		tz = "Asia/Shanghai"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, o.Username, o.Password, o.Database, sslMode, tz)
}

// NewPostgres 使用给定的选项创建一个新的 PostgreSQL gorm 数据库实例.
func NewPostgres(opts *PostgresOptions) (*gorm.DB, error) {
	logLevel := logger.Silent
	if opts.LogLevel != 0 {
		logLevel = logger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(postgres.Open(opts.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
