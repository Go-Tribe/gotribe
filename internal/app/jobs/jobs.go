// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

func InitCron() {
	secondParser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour |
			cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	Cron = cron.New(cron.WithParser(secondParser))
	JobRun(Cron)
}

func JobRun(job *cron.Cron) {
	job.AddFunc("@every 10s", func() {
		//exampleJob()
	})
}
