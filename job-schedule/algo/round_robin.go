package algo

import job "github.com/gagraler/pkg/job-schedule"

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 0:57
 * @file: round_robin.go
 * @description: 轮转调度算法
 */

func RoundRobin(jobs []*job.Job, quantum int) {
	time := 0
	queue := make([]*job.Job, len(jobs))
	copy(queue, jobs)

	for len(queue) > 0 {
		j := queue[0]
		queue = queue[1:]
		if time < j.Arrival {
			time = j.Arrival
		}

		executionTime := min(quantum, j.Burst)
		j.Burst -= executionTime
		j.StartTime = time
		j.WaitTime = time - j.Arrival
		time += executionTime
		j.Arrival = time
		if j.Burst > 0 {
			queue = append(queue, j)
		} else {
			j.FinishTime = time
		}
	}
}
