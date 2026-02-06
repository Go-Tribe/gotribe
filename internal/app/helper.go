// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package app

import (
	"os"
	"path/filepath"
	"strings"

	"gotribe/internal/app/store"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/db"
	"gotribe/pkg/upload"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	// recommendedHomeDir 定义放置 app 服务配置的默认目录.
	recommendedHomeDir = ".configs"

	// defaultConfigName 默认配置文件名（不含扩展名），viper 会按顺序尝试 .yaml / .yml 等
	defaultConfigName = "config"
)

// initConfig 设置需要读取的配置文件名、环境变量，并读取配置文件内容到 viper 中.
// 未指定 -c 时，按以下顺序查找 config.yaml/config.yml：当前目录、可执行文件所在目录、可执行文件上级目录、$HOME/.configs
func initConfig() {
	if cfgFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(cfgFile)
	} else {
		// 1. 当前工作目录（开发时在项目根执行 make run 或 ./gotribe 时生效）
		viper.AddConfigPath(".")

		// 2. 可执行文件所在目录（部署时二进制与 config 同目录）
		if execPath, err := os.Executable(); err == nil {
			viper.AddConfigPath(filepath.Dir(execPath))
			// 3. 可执行文件上级目录（部署时如二进制在 bin/ 子目录，config 在项目根）
			viper.AddConfigPath(filepath.Join(filepath.Dir(execPath), ".."))
		}

		// 4. 用户主目录下的 .configs（统一放置多项目配置）
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为 GoTribe，如果是 app，将自动转变为大写。
	viper.SetEnvPrefix("GoTribe")

	// 以下 2 行，将 viper.Get(key) key 字符串中 '.' 和 '-' 替换为 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	log.Debugw("Using config file", "file", viper.ConfigFileUsed())

	// 配置文件热加载：监听变更后全量重新读取
	if viper.ConfigFileUsed() != "" {
		viper.OnConfigChange(func(e fsnotify.Event) {
			if e.Op != fsnotify.Create && e.Op != fsnotify.Write {
				return
			}
			if err := viper.ReadInConfig(); err != nil {
				log.Errorw("Failed to reload config file", "file", e.Name, "err", err)
				return
			}
			log.Infow("Config file reloaded", "file", e.Name)
		})
		viper.WatchConfig()
	}
}

// logOptions 从 viper 中读取日志配置，构建 `*log.Options` 并返回.
// 注意：`viper.Get<Type>()` 中 key 的名字需要使用 `.` 分割，以跟 YAML 中保持相同的缩进.
func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

// initStore 读取 db 配置，根据 db.type 创建 MySQL 或 PostgreSQL 的 gorm.DB 实例，并初始化 app store 层.
func initStore() error {
	dbType := viper.GetString("db.type")
	if dbType == "" {
		dbType = "mysql"
	}

	var ins *gorm.DB
	var err error

	switch dbType {
	case "postgres", "pg", "postgresql":
		pgOptions := &db.PostgresOptions{
			Host:                  viper.GetString("db.host"),
			Port:                  viper.GetInt("db.port"),
			Username:              viper.GetString("db.username"),
			Password:              viper.GetString("db.password"),
			Database:              viper.GetString("db.database"),
			SSLMode:               viper.GetString("db.sslmode"),
			Timezone:              viper.GetString("db.timezone"),
			MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
			MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
			MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
			LogLevel:              viper.GetInt("db.log-level"),
		}
		ins, err = db.NewPostgres(pgOptions)
	default:
		mysqlOptions := &db.MySQLOptions{
			Host:                  viper.GetString("db.host"),
			Username:              viper.GetString("db.username"),
			Password:              viper.GetString("db.password"),
			Database:              viper.GetString("db.database"),
			MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
			MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
			MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
			LogLevel:              viper.GetInt("db.log-level"),
		}
		ins, err = db.NewMySQL(mysqlOptions)
	}

	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}

// initUpload 根据配置初始化全局上传服务，供 controller 复用
func initUpload() {
	provider := viper.GetString("upload-file.provider")
	if provider == "" {
		if viper.GetBool("enable-oss") {
			provider = "oss"
		} else {
			provider = "qiniu"
		}
	}
	opts := &upload.Options{
		Provider:  upload.Provider(provider),
		Endpoint:  viper.GetString("upload-file.endpoint"),
		AccessKey: viper.GetString("upload-file.access-key"),
		SecretKey: viper.GetString("upload-file.secret-key"),
		Bucket:    viper.GetString("upload-file.bucket"),
		Region:    viper.GetString("upload-file.region"),
	}
	svc, err := upload.NewService(opts)
	if err != nil {
		log.Warnw("Upload service not initialized, upload API will create per-request", "err", err)
		return
	}
	upload.SetDefaultService(svc)
}
