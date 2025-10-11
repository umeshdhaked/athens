package cron

import (
	"github.com/umeshdhaked/athens/pkg/logger"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type job struct {
	scheduler gocron.Scheduler
}

func (j *job) Initialize(interval time.Duration, startDelay time.Duration, executor IExecutor) {

	durationJob := gocron.DurationJob(interval)

	startAt := gocron.WithStartImmediately()
	if startDelay != 0 {
		startAt = gocron.WithStartDateTime(time.Now().Add(startDelay))
	}

	_, err := j.scheduler.NewJob(
		durationJob,
		gocron.NewTask(executor.JobExecutor),
		gocron.WithStartAt(startAt),
		gocron.WithSingletonMode(gocron.LimitModeReschedule))
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("err: ")
	}

	j.scheduler.Start()
}
