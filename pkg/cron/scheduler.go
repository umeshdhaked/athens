package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct{}

func (sh *Scheduler) NewScheduler() ICronJob {
	// create a scheduler
	s, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		log.Fatal("failed initialising cron job")
	}

	return &job{
		scheduler: s,
	}
}
