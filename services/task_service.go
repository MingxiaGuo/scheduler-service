package services

import (
	"scheduler-service/dto"
	"scheduler-service/models"
	"scheduler-service/scheduler"
	"sync"

	"github.com/google/uuid"
)

type TaskService struct {
	mu sync.RWMutex
	// tasks            []*models.Task
	completedTasks   []*models.Task
	scheduleHistory  []models.ScheduleResult
	schedulerManager *scheduler.SchedulerManager
	bandwidth        int
	currentTime      int
	isRunning        bool
}

func NewTaskService(bandwidth int) *TaskService {
	return &TaskService{
		// tasks:            make([]*models.Task, 0),
		completedTasks:   make([]*models.Task, 0),
		scheduleHistory:  make([]models.ScheduleResult, 0),
		schedulerManager: scheduler.NewSchedulerManager(),
		bandwidth:        bandwidth,
		currentTime:      0,
		isRunning:        false,
	}
}

func (ts *TaskService) SubmitTasks(timeSlices []int) (*dto.TaskSubmissionResponse, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// ts.schedulerManager.GetCurrentScheduler().AddTasks(timeSlices)

	// // startIndex := len(ts.tasks) + len(ts.completedTasks)
	for _, timeSlice := range timeSlices {
		// task := models.NewTask(startIndex+i, timeSlice)
		task := models.NewTask(timeSlice)
		// ts.tasks = append(ts.tasks, task)
		ts.schedulerManager.GetCurrentScheduler().AddTasks(*task)
	}

	jobID := uuid.New().String()[:8]
	return &dto.TaskSubmissionResponse{
		JobID:     jobID,
		Message:   "Task submitted successfully",
		TaskCount: len(timeSlices),
	}, nil
}

func (ts *TaskService) GetStatus() *dto.StatusResponse {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return &dto.StatusResponse{
		CurrentTime:     ts.currentTime,
		ScheduleHistory: ts.scheduleHistory,
		ActiveTasks:     ts.getActiveTasksCopy(),
		CompletedTasks:  ts.getCompletedTasksCopy(),
		CurrentStrategy: ts.schedulerManager.GetCurrentScheduler().GetName(),
	}
}

func (ts *TaskService) SwitchScheduler(strategy string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	return ts.schedulerManager.SwitchScheduler(strategy)
}

func (ts *TaskService) ExecuteSchedulingCycle() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.schedulerManager.GetCurrentScheduler().GetTasksLen() == 0 {
		return
	}

	scheduler := ts.schedulerManager.GetCurrentScheduler()
	scheduledTasks := scheduler.Schedule(ts.bandwidth)

	if len(scheduledTasks) > 0 {
		var indexes []int
		var remainingTimes []int

		for _, task := range scheduledTasks {
			indexes = append(indexes, task.Index)
			remainingTimes = append(remainingTimes, task.RemainingTime)
		}

		result := models.ScheduleResult{
			Time:           ts.currentTime,
			TaskIndexes:    indexes,
			RemainingTimes: remainingTimes,
		}
		ts.scheduleHistory = append(ts.scheduleHistory, result)
	}

	ts.moveCompletedTasks(scheduledTasks)
	ts.currentTime++
}

func (ts *TaskService) moveCompletedTasks(tasks []*models.Task) {
	for _, task := range tasks {
		if task.IsCompleted {
			ts.completedTasks = append(ts.completedTasks, task)
		}
	}
}

func (ts *TaskService) getActiveTasksCopy() []models.Task {
	var result []models.Task
	scheduler := ts.schedulerManager.GetCurrentScheduler()
	var tempTasks []models.Task
	for scheduler.GetTasksLen() > 0 {
		task, exists := scheduler.GetNextTask()
		if exists && !task.IsCompleted {
			tempTasks = append(tempTasks, task)
			result = append(result, task)
		}
	}
	for _, task := range tempTasks {
		scheduler.AddTasks(task)
	}

	return result
}

func (ts *TaskService) getCompletedTasksCopy() []models.Task {
	var result []models.Task
	for _, task := range ts.completedTasks {
		result = append(result, *task)
	}
	return result
}

func (ts *TaskService) HasActiveTasks() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.schedulerManager.GetCurrentScheduler().GetTasksLen() > 0
}

func (ts *TaskService) GetAvailableStrategies() []string {
	return ts.schedulerManager.GetAvailableStrategies()
}
