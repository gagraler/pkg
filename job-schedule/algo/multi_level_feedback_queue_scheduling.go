package algo

import (
	job "github.com/gagraler/pkg/job-schedule"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 1:03
 * @file: multi_level_feedback_queue_scheduling.go
 * @description: 多级反馈队列调度算法
 */

func MultilevelFeedbackQueueScheduling(jobs []job.Job) {
	quantum := 2
	time := 0
	queues := [][]job.Job{jobs, {}, {}}
	for len(queues[0]) > 0 || len(queues[1]) > 0 || len(queues[2]) > 0 {
		for i := 0; i < len(queues); i++ {
			if len(queues[i]) == 0 {
				continue
			}
			j := queues[i][0]
			queues[i] = queues[i][1:]
			if time < j.Arrival {
				time = j.Arrival
			}
			executionTime := min(quantum*(i+1), j.Burst)
			j.Burst -= executionTime
			j.WaitTime = time - j.Arrival
			j.StartTime = time
			time += executionTime
			j.Arrival = time
			if j.Burst > 0 {
				if i < len(queues)-1 {
					queues[i+1] = append(queues[i+1], j)
				} else {
					queues[i] = append(queues[i], j)
				}
			} else {
				j.FinishTime = time
			}
		}
	}
}
