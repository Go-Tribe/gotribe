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

// LLMItem 配置中的单个 LLM，仅含 name，实际 api_key/base_url 等从 llm_model 表按 name 查询.
type LLMItem struct {
	Name string `mapstructure:"name" json:"name"`
}

// ModelGroup 配置中的一层：model、strategy、llms（仅 name）.
type ModelGroup struct {
	Model    string    `mapstructure:"model" json:"model"`
	Strategy string    `mapstructure:"strategy" json:"strategy"`
	LLMs     []LLMItem `mapstructure:"llms" json:"llms"`
}

// Config 对应配置最外层的 llms，包含 models 数组.
type Config struct {
	Models []ModelGroup `mapstructure:"models" json:"models"`
}

var (
	ErrLLMConfigNotFound = errors.New("llm config not found")
	ErrNoBackend         = errors.New("no llm in model group")
)

// GetConfig 从 viper 读取 llms 配置（key 为 llms）.
func GetConfig() (*Config, error) {
	var cfg Config
	if err := viper.UnmarshalKey("llms", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SelectBackendName 根据请求的 model（对应配置中的 model）和 strategy 选出一个后端名称.
// 返回的 name 用于在 llm_model 表中查询具体配置（取第一条）.
func SelectBackendName(modelID string) (string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return "", err
	}
	for i := range cfg.Models {
		g := &cfg.Models[i]
		if g.Model == modelID {
			if len(g.LLMs) == 0 {
				return "", ErrNoBackend
			}
			var name string
			switch g.Strategy {
			case strategyRandom:
				name = g.LLMs[rand.Intn(len(g.LLMs))].Name
			default:
				name = g.LLMs[0].Name
			}
			return name, nil
		}
	}
	return "", ErrLLMConfigNotFound
}

// GetMinPointsToChat 返回配置的“少于多少点数不可对话”的阈值，未配置或<=0 视为 0.
func GetMinPointsToChat() float64 {
	v := viper.GetFloat64("llms.min_points_to_chat")
	if v <= 0 {
		return 0
	}
	return v
}

// GetWhiteUsernames 返回配置的白名单用户名列表，这些用户不扣积分，扣费类型记为白名单用户.
func GetWhiteUsernames() []string {
	v := viper.GetStringSlice("llms.white_usernames")
	if v == nil {
		return nil
	}
	return v
}
