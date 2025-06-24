package scheduler

import (
	"testing"
)

func TestNewSchedulerManager(t *testing.T) {
	manager := NewSchedulerManager()

	// 验证默认调度器是FIFO
	if manager.GetCurrentScheduler().GetName() != "FIFO" {
		t.Errorf("Expected default scheduler to be FIFO, got %s",
			manager.GetCurrentScheduler().GetName())
	}

	// 验证可用的调度策略
	strategies := manager.GetAvailableStrategies()
	expectedStrategies := map[string]bool{"FIFO": true, "SRTF": true}

	if len(strategies) != len(expectedStrategies) {
		t.Errorf("Expected %d strategies, got %d",
			len(expectedStrategies), len(strategies))
	}

	for _, strategy := range strategies {
		if !expectedStrategies[strategy] {
			t.Errorf("Unexpected strategy: %s", strategy)
		}
	}
}

func TestSchedulerManager_SwitchScheduler(t *testing.T) {
	manager := NewSchedulerManager()

	// 切换到SRTF
	err := manager.SwitchScheduler("SRTF")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if manager.GetCurrentScheduler().GetName() != "SRTF" {
		t.Errorf("Expected scheduler to be SRTF, got %s",
			manager.GetCurrentScheduler().GetName())
	}

	// 切换到无效的调度器
	err = manager.SwitchScheduler("INVALID")
	if err == nil {
		t.Error("Expected error for invalid scheduler, got nil")
	}

	// 确保调度器没有改变
	if manager.GetCurrentScheduler().GetName() != "SRTF" {
		t.Errorf("Expected scheduler to remain SRTF, got %s",
			manager.GetCurrentScheduler().GetName())
	}
}
