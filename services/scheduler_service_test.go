package services

import (
	"testing"
	"time"
)

func TestNewSchedulerService(t *testing.T) {
	taskService := NewTaskService(5)
	schedulerService := NewSchedulerService(taskService)

	if schedulerService == nil {
		t.Fatal("Expected non-nil SchedulerService")
	}

	if schedulerService.taskService != taskService {
		t.Error("TaskService not properly set")
	}

	if schedulerService.isRunning {
		t.Error("Expected isRunning to be false initially")
	}

	if schedulerService.stopChan == nil {
		t.Error("Expected stopChan to be initialized")
	}
}

func TestSchedulerService_StartStop(t *testing.T) {
	taskService := NewTaskService(5)
	schedulerService := NewSchedulerService(taskService)

	// 测试启动
	schedulerService.Start()
	if !schedulerService.isRunning {
		t.Error("Expected isRunning to be true after Start()")
	}

	// 测试重复启动
	schedulerService.Start() // 不应该有任何影响
	if !schedulerService.isRunning {
		t.Error("Expected isRunning to remain true after second Start()")
	}

	// 测试停止
	schedulerService.Stop()
	if schedulerService.isRunning {
		t.Error("Expected isRunning to be false after Stop()")
	}

	// 测试重复停止
	schedulerService.Stop() // 不应该有任何影响
	if schedulerService.isRunning {
		t.Error("Expected isRunning to remain false after second Stop()")
	}
}

func TestSchedulerService_Run(t *testing.T) {
	taskService := NewTaskService(5)

	// 提交任务
	timeSlices := []int{3, 5, 2}
	_, err := taskService.SubmitTasks(timeSlices)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	schedulerService := NewSchedulerService(taskService)

	// 启动调度器服务
	schedulerService.Start()

	// 等待足够的时间让调度器执行几个周期
	time.Sleep(3 * time.Second)

	// 停止调度器服务
	schedulerService.Stop()

	// 验证任务状态
	status := taskService.GetStatus()

	// 验证当前时间大于0（表示调度器已经运行）
	if status.CurrentTime <= 0 {
		t.Errorf("Expected CurrentTime to be greater than 0, got %d", status.CurrentTime)
	}

	// 验证调度历史记录不为空
	if len(status.ScheduleHistory) == 0 {
		t.Error("Expected non-empty schedule history")
	}

	// 验证任务已经被处理（要么完成，要么活动）
	totalTasks := len(status.ActiveTasks) + len(status.CompletedTasks)
	if totalTasks != len(timeSlices) {
		t.Errorf("Expected %d total tasks, got %d", len(timeSlices), totalTasks)
	}
}
