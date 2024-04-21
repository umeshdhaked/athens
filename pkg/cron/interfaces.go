package cron

import "time"

type ICronJob interface {
	Initialize(interval time.Duration, startDelay time.Duration, executor IExecutor)
}

type ICronScheduler interface {
	NewScheduler() ICronJob
}

type IExecutor interface {
	JobExecutor()
}
