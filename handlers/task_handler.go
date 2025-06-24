package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler-service/dto"
	"scheduler-service/services"
	"scheduler-service/utils"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (th *TaskHandler) SubmitTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var timeSlices []int
	if err := json.NewDecoder(r.Body).Decode(&timeSlices); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request format, expected integer array")
		return
	}

	if len(timeSlices) == 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Task list cannot be empty")
		return
	}

	response, err := th.taskService.SubmitTasks(timeSlices)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to submit tasks")
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (th *TaskHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	status := th.taskService.GetStatus()
	utils.WriteJSONResponse(w, http.StatusOK, status)
}

func (th *TaskHandler) SwitchScheduler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var req dto.SchedulerSwitchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Strategy == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Scheduler strategy cannot be empty")
		return
	}

	availableStrategies := th.taskService.GetAvailableStrategies()
	validStrategy := false
	for _, strategy := range availableStrategies {
		if strategy == req.Strategy {
			validStrategy = true
			break
		}
	}

	if !validStrategy {
		utils.WriteErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unsupported scheduler strategy: %s, available strategies: %v", req.Strategy, availableStrategies))
		return
	}

	if err := th.taskService.SwitchScheduler(req.Strategy); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to switch scheduler strategy")
		return
	}

	response := map[string]interface{}{
		"message":          fmt.Sprintf("Scheduler strategy switched to: %s", req.Strategy),
		"current_strategy": req.Strategy,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}
