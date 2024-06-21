package cron

import (
	"github.com/gagraler/pkg/logger"
	"github.com/go-co-op/gocron/v2"
	"time"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/16 19:02
 * @file: cron.go
 * @description: cron
 */

var log = logger.SugaredLogger()

type Cron struct {
	Expr     string
	TaskFunc func()
}

func (c *Cron) New() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Debugf("cron.NewScheduler err : %v", err)
		return
	}

	j, err := s.NewJob(
		gocron.CronJob(c.Expr, true),
		gocron.NewTask(
			c.TaskFunc,
		))
	if err != nil {
		log.Debugf("NewJob err : %v", err)
		return
	}
	log.Infof("job id: %s, name: %s", j.ID(), j.Name())

	s.Start()
	select {
	case <-time.After(time.Minute):
	}

	err = s.Shutdown()
	if err != nil {
		return
	}
}
