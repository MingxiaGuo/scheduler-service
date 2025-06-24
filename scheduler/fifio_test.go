package scheduler

import (
	"scheduler-service/models"
	"testing"
	"time"
)

func TestFIFOScheduler_GetName(t *testing.T) {
	fifo := NewFIFOScheduler()
	if fifo.GetName() != "FIFO" {
		t.Errorf("Expected scheduler name to be FIFO, got %s", fifo.GetName())
	}
}

func TestFIFOScheduler_Schedule(t *testing.T) {
	tests := []struct {
		name            string
		tasks           []models.Task
		bandwidth       int
		expectedIndexes []int
		expectedRemains []int
	}{
		{
			name: "Single task under bandwidth",
			tasks: []models.Task{
				{Index: 0, RemainingTime: 3, IsCompleted: false, CreatedTime: time.Now()},
			},
			bandwidth:       5,
			expectedIndexes: []int{0},
			expectedRemains: []int{0},
		},
		{
			name: "Multiple tasks under bandwidth",
			tasks: []models.Task{
				{Index: 0, RemainingTime: 2, IsCompleted: false, CreatedTime: time.Now()},
				{Index: 1, RemainingTime: 3, IsCompleted: false, CreatedTime: time.Now().Add(time.Millisecond)},
			},
			bandwidth:       5,
			expectedIndexes: []int{0, 1},
			expectedRemains: []int{0, 0},
		},
		{
			name: "Multiple tasks over bandwidth",
			tasks: []models.Task{
				{Index: 0, RemainingTime: 3, IsCompleted: false, CreatedTime: time.Now()},
				{Index: 1, RemainingTime: 4, IsCompleted: false, CreatedTime: time.Now().Add(time.Millisecond)},
				{Index: 2, RemainingTime: 2, IsCompleted: false, CreatedTime: time.Now().Add(2 * time.Millisecond)},
			},
			bandwidth:       5,
			expectedIndexes: []int{0, 1},
			expectedRemains: []int{0, 2},
		},
		{
			name: "Skip completed tasks",
			tasks: []models.Task{
				{Index: 0, RemainingTime: 0, IsCompleted: true, CreatedTime: time.Now()},
				{Index: 1, RemainingTime: 3, IsCompleted: false, CreatedTime: time.Now().Add(time.Millisecond)},
			},
			bandwidth:       5,
			expectedIndexes: []int{1},
			expectedRemains: []int{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create new scheduler for each test
			fifo := NewFIFOScheduler()

			// Add tasks to scheduler
			for _, task := range tc.tasks {
				fifo.AddTasks(task)
			}

			// Execute scheduling
			scheduledTasks := fifo.Schedule(tc.bandwidth)

			if len(scheduledTasks) != len(tc.expectedIndexes) {
				t.Fatalf("Expected %d scheduled tasks, got %d",
					len(tc.expectedIndexes), len(scheduledTasks))
			}

			for i, task := range scheduledTasks {
				if task.Index != tc.expectedIndexes[i] {
					t.Errorf("Expected task index %d, got %d",
						tc.expectedIndexes[i], task.Index)
				}
				if task.RemainingTime != tc.expectedRemains[i] {
					t.Errorf("Expected remaining time %d, got %d",
						tc.expectedRemains[i], task.RemainingTime)
				}
			}
		})
	}
}
