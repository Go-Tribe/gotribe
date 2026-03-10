// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// Options 发件配置（SMTP）.
type Options struct {
	Host     string `json:"host"`     // SMTP 主机，如 smtp.qq.com
	Port     int    `json:"port"`     // 端口，如 587、465
	From     string `json:"from"`     // 发件人邮箱
	FromName string `json:"fromName"` // 发件人显示名称（如 "GoTribe"），空则仅显示邮箱
	Password string `json:"password"` // 授权码或密码
	UseTLS   bool   `json:"useTLS"`   // 是否使用 TLS（465 一般为 true）
}

// Send 使用给定配置发送一封邮件.
// to 收件人邮箱，subject 主题，body 正文（纯文本）.
func Send(opts *Options, to, subject, body string) error {
	if opts == nil || opts.Host == "" || opts.From == "" {
		return fmt.Errorf("email options incomplete")
	}
	port := opts.Port
	if port == 0 {
		port = 587
	}
	addr := fmt.Sprintf("%s:%d", opts.Host, port)
	auth := smtp.PlainAuth("", opts.From, opts.Password, opts.Host)
	fromHeader := formatAddr(opts.FromName, opts.From)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", fromHeader, to, subject, body))
	// 端口 465 为 SMTPS，必须先 TLS 直连再发 SMTP（与 Python SMTP_SSL 一致）；标准库 SendMail 不会先握手 TLS，会导致 EOF
	if port == 465 {
		return sendMailSMTPS(addr, opts, auth, to, msg)
	}
	return smtp.SendMail(addr, auth, opts.From, []string{to}, msg)
}

// sendMailSMTPS 用于 465 端口 SMTPS：先 TLS 连接再发 SMTP（与 Python SMTP_SSL 行为一致）.
func sendMailSMTPS(addr string, opts *Options, auth smtp.Auth, to string, msg []byte) error {
	tlsConfig := &tls.Config{ServerName: opts.Host}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()
	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(opts.From); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

// formatAddr 格式化为 "显示名" <邮箱>，显示名为空则只返回邮箱.
func formatAddr(name, addr string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return addr
	}
	// 名称中含双引号或反斜杠需转义
	name = strings.ReplaceAll(name, "\\", "\\\\")
	name = strings.ReplaceAll(name, "\"", "\\\"")
	return fmt.Sprintf("\"%s\" <%s>", name, addr)
}

// BuildVerificationMailBody 生成验证码邮件正文（可复用）.
func BuildVerificationMailBody(code string, expireMinutes int) string {
	return fmt.Sprintf("您的验证码是：%s，%d 分钟内有效。如非本人操作请忽略。", code, expireMinutes)
}

// DefaultSubject 验证码邮件默认主题
const DefaultSubject = "验证码"
