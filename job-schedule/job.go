package job_schedule

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/6 0:46
 * @file: job.go
 * @description: job
 */

type Job struct {
	ID         int // id
	Arrival    int // 到达时间
	Burst      int // 运行时间
	Priority   int // 优先级
	WaitTime   int // 等待时间
	StartTime  int // 开始时间
	FinishTime int // 完成时间
}
