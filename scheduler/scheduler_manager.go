package scheduler

import (
	"fmt"
	"scheduler-service/models"
)

type SchedulerManager struct {
	schedulers map[string]Scheduler
	current    Scheduler
}

func NewSchedulerManager() *SchedulerManager {
	fifoScheduler := NewFIFOScheduler()
	srtfScheduler := NewSRTFScheduler()

	schedulers := map[string]Scheduler{
		"FIFO": fifoScheduler,
		"SRTF": srtfScheduler,
	}

	return &SchedulerManager{
		schedulers: schedulers,
		current:    fifoScheduler, // 默认使用FIFO
	}
}

func (sm *SchedulerManager) SwitchScheduler(strategy string) error {
	newScheduler, exists := sm.schedulers[strategy]
	if !exists {
		return fmt.Errorf("unsupported scheduler strategy: %s", strategy)
	}

	if sm.current.GetName() != newScheduler.GetName() {
		sm.migrateTasks(sm.current, newScheduler)
	}
	sm.current = newScheduler

	return nil
}

func (sm *SchedulerManager) migrateTasks(oldScheduler, newScheduler Scheduler) {
	// 从旧调度器中提取所有未完成的任务
	var tasks []models.Task
	for {
		task, exists := oldScheduler.GetNextTask()
		if !exists {
			break
		}
		if !task.IsCompleted {
			tasks = append(tasks, task)
		}
	}

	// 将任务添加到新调度器
	for _, task := range tasks {
		newScheduler.AddTasks(task)
	}
}

func (sm *SchedulerManager) GetCurrentScheduler() Scheduler {
	return sm.current
}

func (sm *SchedulerManager) GetAvailableStrategies() []string {
	var strategies []string
	for name := range sm.schedulers {
		strategies = append(strategies, name)
	}
	return strategies
}
