// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package main

import (
	"os"

	"gotribe/internal/gotribe"

	_ "go.uber.org/automaxprocs"
)

// Go 程序的默认入口函数(主函数).
func main() {
	command := gotribe.NewGoTribeCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
