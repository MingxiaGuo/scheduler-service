package services

import (
	"fmt"
	"sync"
	"time"
)

type SchedulerService struct {
	taskService *TaskService
	stopChan    chan bool
	isRunning   bool
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

func NewSchedulerService(taskService *TaskService) *SchedulerService {
	return &SchedulerService{
		taskService: taskService,
		stopChan:    make(chan bool),
		isRunning:   false,
	}
}

func (ss *SchedulerService) Start() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if ss.isRunning {
		return
	}

	ss.isRunning = true
	ss.wg.Add(1)
	go ss.run()
	fmt.Println("Scheduler service started")
}

func (ss *SchedulerService) Stop() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.isRunning {
		return
	}

	ss.stopChan <- true
	ss.isRunning = false
	fmt.Println("Scheduler service stopped")
}

func (ss *SchedulerService) run() {
	defer ss.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ss.taskService.HasActiveTasks() {
				ss.taskService.ExecuteSchedulingCycle()
				ss.printCurrentStatus()
			}
		case <-ss.stopChan:
			return
		}
	}
}

func (ss *SchedulerService) printCurrentStatus() {
	status := ss.taskService.GetStatus()
	if len(status.ScheduleHistory) == 0 {
		return
	}

	latestResult := status.ScheduleHistory[len(status.ScheduleHistory)-1]
	fmt.Printf("Time: %d, Executed Task Indexes: %v, Remaining Times: %v\n",
		latestResult.Time,
		latestResult.TaskIndexes,
		latestResult.RemainingTimes)
}

func (ss *SchedulerService) GracefulStop() {
	ss.mu.Lock()
	if !ss.isRunning {
		ss.mu.Unlock()
		return
	}

	fmt.Println("Gracefully stopping scheduler service...")
	ss.stopChan <- true
	ss.isRunning = false
	ss.mu.Unlock()

	ss.wg.Wait()

	ss.waitForCurrentTasks()

	fmt.Println("Scheduler service gracefully stopped")
}

func (ss *SchedulerService) waitForCurrentTasks() {
	timeout := time.After(10 * time.Second) // 最多等待10秒
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("Timeout waiting for tasks to complete, forcing shutdown")
			return
		case <-ticker.C:
			if !ss.taskService.HasActiveTasks() {
				fmt.Println("All tasks completed")
				return
			}
		}
	}
}
