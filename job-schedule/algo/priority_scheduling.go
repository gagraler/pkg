package algo

import (
	job "github.com/gagraler/pkg/job-schedule"
	"sort"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 1:00
 * @file: priority_scheduling.go
 * @description: 优先级调度算法
 */

func PriorityScheduling(jobs []job.Job) {
	time := 0
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Priority < jobs[j].Priority
	})
	for _, j := range jobs {
		if time < j.Arrival {
			time = j.Arrival
		}
		j.WaitTime = time - j.Arrival
		j.StartTime = time
		time += j.Burst
		j.FinishTime = time
	}
}
