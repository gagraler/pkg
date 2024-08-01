package scheduler

import (
	"container/heap"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/7/14 20:47
 * @file: scheduler.go
 * @description: 任务调度器
 */

type TaskStatus string

const (
	Pending   TaskStatus = "pending"
	Running   TaskStatus = "running"
	Completed TaskStatus = "completed"
	NotFound  TaskStatus = "not found"
)

// Task 任务
type Task struct {
	Id        string
	Priority  int
	ExecuteAt time.Time
	Task      func()
	Status    TaskStatus
}

// TaskQueue 任务优先队列
type TaskQueue []*Task

// Len 返回队列中的元素数量
func (tq *TaskQueue) Len() int { return len(*tq) }

// Less 如果索引i处的任务比索引j处的任务优先级高，返回true
func (tq *TaskQueue) Less(i, j int) bool {
	if (*tq)[i].ExecuteAt.Equal((*tq)[j].ExecuteAt) {
		return (*tq)[i].Priority > (*tq)[j].Priority
	}
	return (*tq)[i].ExecuteAt.Before((*tq)[j].ExecuteAt)
}

// Swap 交换给定索引处的元素
func (tq *TaskQueue) Swap(i, j int) { (*tq)[i], (*tq)[j] = (*tq)[j], (*tq)[i] }

// Push 向队列中添加一个元素
func (tq *TaskQueue) Push(x interface{}) {
	*tq = append(*tq, x.(*Task))
}

// Pop 移除并返回优先级最高的元素
func (tq *TaskQueue) Pop() interface{} {
	old := *tq
	n := len(old)
	task := old[n-1]
	*tq = old[0 : n-1]
	return task
}

// Scheduler 调度器
type Scheduler struct {
	taskQueue *TaskQueue
	stopChan  chan bool
	mu        sync.Mutex
}

// NewScheduler 创建一个新的调度器
func NewScheduler() *Scheduler {
	tq := &TaskQueue{}
	heap.Init(tq)
	return &Scheduler{
		taskQueue: tq,
		stopChan:  make(chan bool),
	}
}

// AddTask adds a new task to the scheduler.
func (s *Scheduler) AddTask(executeAt time.Time, priority int, task func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap.Push(s.taskQueue, &Task{
		Id:        getId(),
		Priority:  priority,
		ExecuteAt: executeAt,
		Task:      task,
		Status:    Pending,
	})
}

// RemoveTask removes a task from the scheduler.
func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range *s.taskQueue {
		if task.Id == id {
			heap.Remove(s.taskQueue, i)
			break
		}
	}
}

// Start starts the scheduler to run tasks at their scheduled times.
func (s *Scheduler) Start() {
	go func() {
		for {
			// 检查是否有任务
			s.mu.Lock()
			if s.taskQueue.Len() == 0 {
				s.mu.Unlock()
				select {
				case <-s.stopChan:
					return
				default:
					time.Sleep(1 * time.Second)
					continue
				}
			}

			// 获取下一个任务
			nextTask := heap.Pop(s.taskQueue).(*Task)
			nextTask.Status = Running
			s.mu.Unlock()

			timeUntilExecute := nextTask.ExecuteAt.Sub(time.Now())
			if timeUntilExecute > 0 {
				select {
				case <-time.After(timeUntilExecute):
				case <-s.stopChan:
					return
				}
			}

			nextTask.Task()
			s.mu.Lock()
			nextTask.Status = Completed
			s.mu.Unlock()
		}
	}()
}

// Stop stops the scheduler from running tasks.
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// GetTaskStatus returns the status of a task.
func (s *Scheduler) GetTaskStatus(id string) TaskStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, task := range *s.taskQueue {
		if task.Id == id {
			return task.Status
		}
	}
	return NotFound
}

// getId generates a new UUID not horizontal line
func getId() string {
	u := uuid.New().String()
	return strings.Replace(u, "-", "", -1)
}
