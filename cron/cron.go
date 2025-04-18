package cron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/16 19:02
 * @file: cron.go
 * @description: cron
 */

type Cron struct {
	Expr     string
	TaskFunc func()
}

func (c *Cron) New() error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("cron.NewScheduler err: %v", err)
	}

	j, err := s.NewJob(
		gocron.CronJob(c.Expr, true),
		gocron.NewTask(
			c.TaskFunc,
		))
	if err != nil {
		return fmt.Errorf("cron.NewJob err: %v", err)
	}
	fmt.Printf("job id: %s, name: %s\n", j.ID(), j.Name())

	s.Start()
	time.Sleep(time.Minute)

	err = s.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
