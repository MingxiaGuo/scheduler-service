package scheduler

import (
	"scheduler-service/models"
)

type SRTFScheduler struct {
	*BaseScheduler
}

func NewSRTFScheduler() *SRTFScheduler {
	heap := make(SRTFTaskHeap, 0)
	baseScheduler := NewBaseScheduler(&heap, "SRTF")
	return &SRTFScheduler{
		BaseScheduler: baseScheduler,
	}
}

type SRTFTaskHeap []models.Task

func (h SRTFTaskHeap) Len() int {
	return len(h)
}

func (h SRTFTaskHeap) Less(i, j int) bool {
	if h[i].RemainingTime == h[j].RemainingTime {
		return h[i].Index < h[j].Index
	}
	return h[i].RemainingTime < h[j].RemainingTime
}

func (h SRTFTaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *SRTFTaskHeap) Push(x interface{}) {
	*h = append(*h, x.(models.Task))
}

func (h *SRTFTaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	*h = old[0 : n-1]
	return task
}
