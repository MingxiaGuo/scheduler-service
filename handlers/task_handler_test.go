package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scheduler-service/dto"

	"scheduler-service/services"
	"testing"
)

func TestTaskHandler_SubmitTasks(t *testing.T) {
	taskService := services.NewTaskService(5)
	taskHandler := NewTaskHandler(taskService)

	tests := []struct {
		name           string
		method         string
		body           []byte
		expectedStatus int
		validateResp   func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid request",
			method:         http.MethodPost,
			body:           []byte(`[3, 5, 2]`),
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var response dto.TaskSubmissionResponse
				err := json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.TaskCount != 3 {
					t.Errorf("Expected task count to be 3, got %d", response.TaskCount)
				}

				if response.JobID == "" {
					t.Error("Expected JobID to be non-empty")
				}
			},
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           []byte(`[3, 5, 2]`),
			expectedStatus: http.StatusMethodNotAllowed,
			validateResp:   nil,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			body:           []byte(`{invalid json}`),
			expectedStatus: http.StatusBadRequest,
			validateResp:   nil,
		},
		{
			name:           "Empty task list",
			method:         http.MethodPost,
			body:           []byte(`[]`),
			expectedStatus: http.StatusBadRequest,
			validateResp:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/tasks", bytes.NewBuffer(tc.body))
			resp := httptest.NewRecorder()

			taskHandler.SubmitTasks(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d",
					tc.expectedStatus, resp.Code)
			}

			if tc.validateResp != nil {
				tc.validateResp(t, resp)
			}
		})
	}
}

func TestTaskHandler_GetStatus(t *testing.T) {
	taskService := services.NewTaskService(5)
	taskHandler := NewTaskHandler(taskService)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/status", nil)
			resp := httptest.NewRecorder()

			taskHandler.GetStatus(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d",
					tc.expectedStatus, resp.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				var status dto.StatusResponse
				err := json.NewDecoder(resp.Body).Decode(&status)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if status.CurrentStrategy != "FIFO" {
					t.Errorf("Expected strategy to be FIFO, got %s",
						status.CurrentStrategy)
				}
			}
		})
	}
}

func TestTaskHandler_SwitchScheduler(t *testing.T) {
	taskService := services.NewTaskService(5)
	taskHandler := NewTaskHandler(taskService)

	tests := []struct {
		name           string
		method         string
		body           []byte
		expectedStatus int
	}{
		{
			name:           "Valid request - SRTF",
			method:         http.MethodPost,
			body:           []byte(`{"strategy":"SRTF"}`),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request - FIFO",
			method:         http.MethodPost,
			body:           []byte(`{"strategy":"FIFO"}`),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           []byte(`{"strategy":"FIFO"}`),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			body:           []byte(`{invalid json}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty strategy",
			method:         http.MethodPost,
			body:           []byte(`{"strategy":""}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid strategy",
			method:         http.MethodPost,
			body:           []byte(`{"strategy":"INVALID"}`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/scheduler", bytes.NewBuffer(tc.body))
			resp := httptest.NewRecorder()

			taskHandler.SwitchScheduler(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d",
					tc.expectedStatus, resp.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				strategy, ok := response["current_strategy"]
				if !ok {
					t.Fatal("Expected current_strategy in response")
				}

				strategySwitchReq := dto.SchedulerSwitchRequest{}
				err = json.Unmarshal(tc.body, &strategySwitchReq)
				if err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}

				if strategy != strategySwitchReq.Strategy {
					t.Errorf("Expected strategy to be %s, got %s",
						strategySwitchReq.Strategy, strategy)
				}
			}
		})
	}
}
