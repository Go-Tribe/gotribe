// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

import "fmt"

// Errno 定义了 app 使用的错误类型.
type Errno struct {
	HTTP    int
	Code    string
	Message map[string]string // 使用 map 存储不同语言的消息
}

// Error 实现 error 接口中的 `Error` 方法.
func (err *Errno) Error() string {
	defaultLang := "zh" // 默认语言
	if err.Message == nil || len(err.Message) == 0 {
		return "Unknown error"
	}
	for lang, msg := range err.Message {
		if lang == defaultLang {
			return msg
		}
	}
	// 如果默认语言不存在，返回第一个可用的语言
	for _, msg := range err.Message {
		return msg
	}
	return "Unknown error"
}

// SetMessage 设置 Errno 类型错误中的 Message 字段.
func (err *Errno) SetMessage(format string, args ...interface{}) *Errno {
	lang := "zh" // 默认语言为中文

	// 检查是否传递了 lang 参数
	if len(args) > 0 && len(format) == 0 {
		lang = format
		format = args[0].(string)
		args = args[1:]
	}

	if err.Message == nil {
		err.Message = make(map[string]string)
	}
	err.Message[lang] = fmt.Sprintf(format, args...)
	return err
}

// Decode 尝试从 err 中解析出业务错误码和错误信息.
func Decode(err error, lang string) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message[lang]
	}

	switch typed := err.(type) {
	case *Errno:
		msg, ok := typed.Message[lang]
		if !ok && typed.Message != nil {
			// 如果指定语言不存在，尝试使用默认语言 "en"
			msg, ok = typed.Message["en"]
			if !ok {
				// 如果默认语言也不存在，返回第一个可用的语言
				for _, m := range typed.Message {
					msg = m
					break
				}
			}
		}
		return typed.HTTP, typed.Code, msg
	default:
		// 对于未知类型的错误，返回 InternalServerError 并保留原始错误信息
		return InternalServerError.HTTP, InternalServerError.Code, err.Error()
	}
}
