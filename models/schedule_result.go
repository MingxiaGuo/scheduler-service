package models

type ScheduleResult struct {
	Time           int   `json:"time"`
	TaskIndexes    []int `json:"task_indexes"`
	RemainingTimes []int `json:"remaining_times"`
}
