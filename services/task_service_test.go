package services

import (
	"testing"
)

func TestNewTaskService(t *testing.T) {
	bandwidth := 5
	service := NewTaskService(bandwidth)

	status := service.GetStatus()
	if status.CurrentTime != 0 {
		t.Errorf("Expected CurrentTime to be 0, got %d", status.CurrentTime)
	}

	if len(status.ActiveTasks) != 0 {
		t.Errorf("Expected 0 active tasks, got %d", len(status.ActiveTasks))
	}

	if len(status.CompletedTasks) != 0 {
		t.Errorf("Expected 0 completed tasks, got %d", len(status.CompletedTasks))
	}

	if status.CurrentStrategy != "FIFO" {
		t.Errorf("Expected default strategy to be FIFO, got %s", status.CurrentStrategy)
	}
}

func TestTaskService_SubmitTasks(t *testing.T) {
	service := NewTaskService(5)

	timeSlices := []int{3, 5, 2}
	response, err := service.SubmitTasks(timeSlices)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.TaskCount != len(timeSlices) {
		t.Errorf("Expected task count to be %d, got %d",
			len(timeSlices), response.TaskCount)
	}

	if response.JobID == "" {
		t.Error("Expected JobID to be non-empty")
	}

	status := service.GetStatus()
	if len(status.ActiveTasks) != len(timeSlices) {
		t.Errorf("Expected %d active tasks, got %d",
			len(timeSlices), len(status.ActiveTasks))
	}
}

func TestTaskService_ExecuteSchedulingCycle(t *testing.T) {
	service := NewTaskService(5)

	// 提交任务
	timeSlices := []int{3, 5, 2}
	_, err := service.SubmitTasks(timeSlices)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// 执行一个调度周期
	service.ExecuteSchedulingCycle()

	status := service.GetStatus()
	if status.CurrentTime != 1 {
		t.Errorf("Expected CurrentTime to be 1, got %d", status.CurrentTime)
	}

	if len(status.ScheduleHistory) != 1 {
		t.Fatalf("Expected 1 schedule history entry, got %d",
			len(status.ScheduleHistory))
	}

	history := status.ScheduleHistory[0]
	if history.Time != 0 {
		t.Errorf("Expected history time to be 0, got %d", history.Time)
	}

	// 验证任务状态
	activeTasks := len(status.ActiveTasks)
	completedTasks := len(status.CompletedTasks)
	expectedActive := 2    // 5带宽，FIFO模式下，第一个任务(3)完成，第二个任务(5)部分完成(剩余3)，第三个任务(2)未处理
	expectedCompleted := 1 // 只有第一个任务完成

	if activeTasks != expectedActive {
		t.Errorf("Expected %d active tasks, got %d",
			expectedActive, activeTasks)
	}

	if completedTasks != expectedCompleted {
		t.Errorf("Expected %d completed tasks, got %d",
			expectedCompleted, completedTasks)
	}
}

func TestTaskService_SwitchScheduler(t *testing.T) {
	service := NewTaskService(5)

	// 切换到SRTF
	err := service.SwitchScheduler("SRTF")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	status := service.GetStatus()
	if status.CurrentStrategy != "SRTF" {
		t.Errorf("Expected strategy to be SRTF, got %s", status.CurrentStrategy)
	}

	// 切换到无效的调度器
	err = service.SwitchScheduler("INVALID")
	if err == nil {
		t.Error("Expected error for invalid scheduler, got nil")
	}

	// 确保调度器没有改变
	status = service.GetStatus()
	if status.CurrentStrategy != "SRTF" {
		t.Errorf("Expected strategy to remain SRTF, got %s", status.CurrentStrategy)
	}
}
