package dto

import "scheduler-service/models"

type TaskRequest struct {
	Tasks    string `json:"tasks"`    // Task List
	Strategy string `json:"strategy"` // Scheduler Strategy: FIFO or SRTF
}

type TaskSubmissionResponse struct {
	JobID     string `json:"job_id"`
	Message   string `json:"message"`
	TaskCount int    `json:"task_count"`
}

type StatusResponse struct {
	CurrentTime     int                     `json:"current_time"`
	ScheduleHistory []models.ScheduleResult `json:"schedule_history"`
	ActiveTasks     []models.Task           `json:"active_tasks"`
	CompletedTasks  []models.Task           `json:"completed_tasks"`
	CurrentStrategy string                  `json:"current_strategy"`
}

type SchedulerSwitchRequest struct {
	Strategy string `json:"strategy" binding:"required"`
}
