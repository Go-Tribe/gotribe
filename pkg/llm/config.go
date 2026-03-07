// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package llm

import (
	"errors"
	"math/rand"

	"github.com/spf13/viper"
)

const strategyRandom = "random"

// Backend 表示一个具体的 LLM 后端（api_key、base_url、model、header）.
// Header 可配置额外请求头，用于部分模型所需的自定义 header.
type Backend struct {
	Type    string            `mapstructure:"type" json:"type"` // openai
	APIKey  string            `mapstructure:"api_key" json:"api_key"`
	BaseURL string            `mapstructure:"base_url" json:"base_url"`
	Model   string            `mapstructure:"model" json:"model"`
	Header  map[string]string `mapstructure:"header" json:"header"`
}

// ModelGroup 表示配置中的一层：id、strategy、以及多个 models.
type ModelGroup struct {
	ID       string    `mapstructure:"id" json:"id"`
	Strategy string    `mapstructure:"strategy" json:"strategy"`
	Models   []Backend `mapstructure:"models" json:"models"`
}

// Config 对应配置最外层的 llms，包含 models 数组.
type Config struct {
	Models []ModelGroup `mapstructure:"models" json:"models"`
}

var (
	ErrLLMConfigNotFound = errors.New("llm config not found")
	ErrNoBackend         = errors.New("no backend in model group")
)

// GetConfig 从 viper 读取 llms 配置（key 为 llms）.
func GetConfig() (*Config, error) {
	var cfg Config
	if err := viper.UnmarshalKey("llms", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SelectBackend 根据请求的 model（对应配置中的 id）和 strategy 选出一个后端.
// 目前仅支持 strategy=random.
func SelectBackend(modelID string) (*Backend, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	for i := range cfg.Models {
		g := &cfg.Models[i]
		if g.ID == modelID {
			if len(g.Models) == 0 {
				return nil, ErrNoBackend
			}
			switch g.Strategy {
			case strategyRandom:
				return &g.Models[rand.Intn(len(g.Models))], nil
			default:
				return &g.Models[0], nil
			}
		}
	}
	return nil, ErrLLMConfigNotFound
}
