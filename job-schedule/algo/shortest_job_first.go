package algo

import (
	job "github.com/gagraler/pkg/job-schedule"
	"sort"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 0:52
 * @file: shortest_job_first.go
 * @description: 最短作业优先
 */

func ShortestJobFirstAlgo(jobs []*job.Job) {
	time := 0
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Burst < jobs[j].Burst
	})

	for _, j := range jobs {
		if time < j.Arrival {
			time = j.Arrival
		}
		j.StartTime = time
		//j.FinishTime = time + j.Burst
		j.WaitTime = time - j.Arrival
		time += j.Burst
		j.FinishTime = time
	}
}
