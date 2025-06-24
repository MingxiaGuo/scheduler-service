package scheduler

import (
	"container/heap"
	"scheduler-service/models"
)

type HeapInterface interface {
	heap.Interface
	Len() int
}

type BaseScheduler struct {
	heap HeapInterface
	name string
}

func NewBaseScheduler(h HeapInterface, name string) *BaseScheduler {
	heap.Init(h)
	return &BaseScheduler{
		heap: h,
		name: name,
	}
}

func (b *BaseScheduler) GetName() string {
	return b.name
}

func (b *BaseScheduler) Schedule(bandwidth int) []*models.Task {
	var scheduledTasks []*models.Task
	usedBandwidth := 0
	var tempTasks []models.Task

	for b.heap.Len() > 0 && usedBandwidth < bandwidth {
		task := heap.Pop(b.heap).(models.Task)
		if task.IsCompleted {
			continue
		}

		availableBandwidth := bandwidth - usedBandwidth
		if availableBandwidth <= 0 {
			tempTasks = append(tempTasks, task)
			break
		}

		allocatedTime := min(task.RemainingTime, availableBandwidth)
		task.Execute(allocatedTime)
		usedBandwidth += allocatedTime

		scheduledTasks = append(scheduledTasks, &task)

		if !task.IsCompleted {
			tempTasks = append(tempTasks, task)
		}
	}

	for _, task := range tempTasks {
		heap.Push(b.heap, task)
	}

	return scheduledTasks
}

func (b *BaseScheduler) AddTasks(task models.Task) {
	if !task.IsCompleted {
		heap.Push(b.heap, task)
	}
}

func (b *BaseScheduler) GetNextTask() (models.Task, bool) {
	for b.heap.Len() > 0 {
		task := heap.Pop(b.heap).(models.Task)
		if !task.IsCompleted {
			return task, true
		}
	}
	return models.Task{}, false
}

func (b *BaseScheduler) GetTasksLen() int {
	return b.heap.Len()
}
