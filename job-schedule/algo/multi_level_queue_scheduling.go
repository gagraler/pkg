package algo

import (
	"fmt"
	job "github.com/gagraler/pkg/job-schedule"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 1:01
 * @file: multi_level_queue_scheduling.go
 * @description: 多级队列调度算法
 */

func MultilevelQueueScheduling(jobs []job.Job) {
	fmt.Println("多级队列调度：")
	// 假设这里有两个队列，一个优先处理短作业，一个优先处理长作业
	var shortQueue []job.Job
	var longQueue []job.Job
	for _, j := range jobs {
		if j.Burst <= 5 {
			shortQueue = append(shortQueue, j)
		} else {
			longQueue = append(longQueue, j)
		}
	}
	// 处理短作业队列
	time := 0
	for _, j := range shortQueue {
		if time < j.Arrival {
			time = j.Arrival
		}
		j.WaitTime = time - j.Arrival
		j.StartTime = time
		time += j.Burst
		j.FinishTime = time
	}

	// 处理长作业队列
	for _, j := range longQueue {
		if time < j.Arrival {
			time = j.Arrival
		}
		j.WaitTime = time - j.Arrival
		j.StartTime = time
		time += j.Burst
		j.FinishTime = time
	}
}
