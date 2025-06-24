package models

import (
	"testing"
)

func TestNewTask(t *testing.T) {
	task := NewTask(10)
	if task.RemainingTime != 10 {
		t.Errorf("Expected RemainingTime to be 10, got %d", task.RemainingTime)
	}
	if task.IsCompleted {
		t.Error("Expected IsCompleted to be false, got true")
	}
}

func TestTaskExecute(t *testing.T) {
	tests := []struct {
		name           string
		initialTime    int
		executeTime    int
		expectedRemain int
		expectedDone   bool
	}{
		{"Partial execution", 10, 4, 6, false},
		{"Exact execution", 5, 5, 0, true},
		{"Over execution", 3, 5, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask(tc.initialTime)
			task.Execute(tc.executeTime)

			if task.RemainingTime != tc.expectedRemain {
				t.Errorf("Expected RemainingTime to be %d, got %d",
					tc.expectedRemain, task.RemainingTime)
			}

			if task.IsCompleted != tc.expectedDone {
				t.Errorf("Expected IsCompleted to be %v, got %v",
					tc.expectedDone, task.IsCompleted)
			}
		})
	}
}
