package scheduler

import "scheduler-service/models"

type Scheduler interface {
	Schedule(bandwidth int) []*models.Task
	GetName() string
	AddTasks(task models.Task)
	GetNextTask() (models.Task, bool)
	GetTasksLen() int
}
