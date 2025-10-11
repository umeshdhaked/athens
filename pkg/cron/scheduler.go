package cron

import (
	"github.com/umeshdhaked/athens/pkg/logger"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct{}

func (sh *Scheduler) NewScheduler() ICronJob {
	// create a scheduler
	s, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("failed initialising cron job")
	}

	return &job{
		scheduler: s,
	}
}
