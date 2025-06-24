// package scheduler

// import (
// 	"scheduler-service/models"

// 	"testing"
// )

// func TestSRTFScheduler_GetName(t *testing.T) {
// 	srtf := NewSRTFScheduler()
// 	if srtf.GetName() != "SRTF" {
// 		t.Errorf("Expected scheduler name to be SRTF, got %s", srtf.GetName())
// 	}
// }

// func TestSRTFScheduler_Schedule(t *testing.T) {
// 	srtf := NewSRTFScheduler()

// 	tests := []struct {
// 		name            string
// 		tasks           []*models.Task
// 		bandwidth       int
// 		expectedIndexes []int
// 		expectedRemains []int
// 	}{
// 		{
// 			name: "Sort by remaining time",
// 			tasks: []*models.Task{
// 				{Index: 0, RemainingTime: 5, IsCompleted: false},
// 				{Index: 1, RemainingTime: 2, IsCompleted: false},
// 				{Index: 2, RemainingTime: 8, IsCompleted: false},
// 			},
// 			bandwidth:       5,
// 			expectedIndexes: []int{1, 0},
// 			expectedRemains: []int{0, 2},
// 		},
// 		{
// 			name: "Equal remaining time, sort by index",
// 			tasks: []*models.Task{
// 				{Index: 1, RemainingTime: 3, IsCompleted: false},
// 				{Index: 0, RemainingTime: 3, IsCompleted: false},
// 			},
// 			bandwidth:       5,
// 			expectedIndexes: []int{0, 1},
// 			expectedRemains: []int{0, 1},
// 		},
// 		{
// 			name: "Skip completed tasks",
// 			tasks: []*models.Task{
// 				{Index: 0, RemainingTime: 0, IsCompleted: true},
// 				{Index: 1, RemainingTime: 3, IsCompleted: false},
// 			},
// 			bandwidth:       5,
// 			expectedIndexes: []int{1},
// 			expectedRemains: []int{0},
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			scheduledTasks := srtf.Schedule(tc.bandwidth)

// 			if len(scheduledTasks) != len(tc.expectedIndexes) {
// 				t.Fatalf("Expected %d scheduled tasks, got %d",
// 					len(tc.expectedIndexes), len(scheduledTasks))
// 			}

//				for i, task := range scheduledTasks {
//					if task.Index != tc.expectedIndexes[i] {
//						t.Errorf("Expected task index %d, got %d",
//							tc.expectedIndexes[i], task.Index)
//					}
//					if task.RemainingTime != tc.expectedRemains[i] {
//						t.Errorf("Expected remaining time %d, got %d",
//							tc.expectedRemains[i], task.RemainingTime)
//					}
//				}
//			})
//		}
//	}
package scheduler

import (
	"scheduler-service/models"
	"testing"
)

func TestSRTFScheduler_GetName(t *testing.T) {
	srtf := NewSRTFScheduler()
	if srtf.GetName() != "SRTF" {
		t.Errorf("Expected scheduler name to be SRTF, got %s", srtf.GetName())
	}
}

func TestSRTFScheduler_Schedule(t *testing.T) {
	tests := []struct {
		name            string
		tasks           []*models.Task
		bandwidth       int
		expectedIndexes []int
		expectedRemains []int
	}{
		{
			name: "Sort by remaining time",
			tasks: []*models.Task{
				{Index: 0, RemainingTime: 5, IsCompleted: false},
				{Index: 1, RemainingTime: 2, IsCompleted: false},
				{Index: 2, RemainingTime: 8, IsCompleted: false},
			},
			bandwidth:       5,
			expectedIndexes: []int{1, 0},
			expectedRemains: []int{0, 2},
		},
		{
			name: "Equal remaining time, sort by index",
			tasks: []*models.Task{
				{Index: 1, RemainingTime: 3, IsCompleted: false},
				{Index: 0, RemainingTime: 3, IsCompleted: false},
			},
			bandwidth:       5,
			expectedIndexes: []int{0, 1},
			expectedRemains: []int{0, 1},
		},
		{
			name: "Skip completed tasks",
			tasks: []*models.Task{
				{Index: 0, RemainingTime: 0, IsCompleted: true},
				{Index: 1, RemainingTime: 3, IsCompleted: false},
			},
			bandwidth:       5,
			expectedIndexes: []int{1},
			expectedRemains: []int{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srtf := NewSRTFScheduler()

			// Add tasks to the scheduler
			for _, task := range tc.tasks {
				srtf.AddTasks(*task)
			}

			scheduledTasks := srtf.Schedule(tc.bandwidth)

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
