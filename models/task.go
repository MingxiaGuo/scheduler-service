package models

import (
	"sync/atomic"
	"time"
)

var globalIndex int32

type Task struct {
	Index         int
	RemainingTime int
	IsCompleted   bool
	CreatedTime   time.Time
}

func NewTask(duration int) *Task {
	idx := atomic.AddInt32(&globalIndex, 1) - 1
	return &Task{
		Index:         int(idx),
		RemainingTime: duration,
		CreatedTime:   time.Now(),
	}
}

func (t *Task) Execute(timeSlice int) {
	if t.RemainingTime > timeSlice {
		t.RemainingTime -= timeSlice
	} else {
		t.RemainingTime = 0
		t.IsCompleted = true
	}
}
