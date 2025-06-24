package scheduler

import (
	"scheduler-service/models"
)

// FIFOScheduler 顺序调度
type FIFOScheduler struct {
	*BaseScheduler
}

func NewFIFOScheduler() *FIFOScheduler {
	heap := make(FIFOTaskHeap, 0)
	baseScheduler := NewBaseScheduler(&heap, "FIFO")
	return &FIFOScheduler{
		BaseScheduler: baseScheduler,
	}
}

type FIFOTaskHeap []models.Task

func (h FIFOTaskHeap) Len() int {
	return len(h)
}

func (h FIFOTaskHeap) Less(i, j int) bool {
	return h[i].CreatedTime.Before(h[j].CreatedTime)
}

func (h FIFOTaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *FIFOTaskHeap) Push(x interface{}) {
	*h = append(*h, x.(models.Task))
}

func (h *FIFOTaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	*h = old[0 : n-1]
	return task
}
