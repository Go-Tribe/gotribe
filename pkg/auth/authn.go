// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package auth

import "golang.org/x/crypto/bcrypt"

// Encrypt 使用 bcrypt 加密纯文本.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)

	return string(hashedBytes), err
}

// Compare 比较密文和明文是否相同.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
