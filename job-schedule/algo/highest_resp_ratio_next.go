package algo

import job "github.com/gagraler/pkg/job-schedule"

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 0:54
 * @file: highest_resp_ratio_next.go
 * @description: 最高响应比优先
 */

func HighestRespRatioNextAlgo(jobs []*job.Job) {
	time := 0
	for len(jobs) > 0 {
		highestResponseRatio := -1.0
		index := 0

		for i, j := range jobs {
			waitTime := time - j.Arrival
			if waitTime < 0 {
				waitTime = 0
			}

			responseRatio := float64(waitTime+j.Burst) / float64(j.Burst)
			if responseRatio > highestResponseRatio {
				highestResponseRatio = responseRatio
				index = i
			}
		}

		job := jobs[index]
		jobs = append(jobs[:index], jobs[index+1:]...)
		if time < job.Arrival {
			time = job.Arrival
		}

		job.StartTime = time
		job.WaitTime = time - job.Arrival
		time += job.Burst
		job.FinishTime = time
	}
}
