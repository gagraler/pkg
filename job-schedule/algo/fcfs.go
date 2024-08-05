package algo

import job "github.com/gagraler/pkg/job-schedule"

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 0:49
 * @file: fcfs.go
 * @description: 先来先服务
 */

func FirstComeToFirstServedAlgo(jobs []*job.Job) {

	time := 0
	for _, j := range jobs {
		if time < j.Arrival {
			time = j.Arrival
		}
		j.StartTime = time
		j.FinishTime = time + j.Burst
		j.WaitTime = time - j.Arrival
		time += j.Burst
		j.FinishTime = time
	}
}
