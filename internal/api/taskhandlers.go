package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gogazub/app/internal/repo"
	"github.com/gogazub/app/internal/service"
)

type TaskResponse struct {
	TaskID string `json:"task_id"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ResultResponse struct {
	Result string `json:"result"`
}

func HandleTaskStatus(w http.ResponseWriter, r *http.Request, service *service.Service) {

	if err := validateToken(r, service); err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	taskID, valid := handleTaskID(w, r)
	if !valid {
		return
	}

	status, err := service.GetTaskService().GetTaskStatus(taskID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	statusStr := mapStatusToString(status)

	response := StatusResponse{Status: statusStr}
	writeJSON(w, http.StatusOK, response)
}

func HandleTaskResult(w http.ResponseWriter, r *http.Request, service *service.Service) {

	if err := validateToken(r, service); err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	taskID, valid := handleTaskID(w, r)
	if !valid {
		return
	}

	result, err := service.GetTaskService().GetTaskResult(taskID)
	if err != nil {
		if strings.Contains(err.Error(), "not ready") {
			writeJSONError(w, http.StatusAccepted, err.Error()) // статус Accepted если результат не готов
		} else {
			writeJSONError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	response := ResultResponse{Result: result}
	writeJSON(w, http.StatusOK, response)
}

func handleTaskID(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return "", false
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		writeJSONError(w, http.StatusBadRequest, "Task ID is required")
		return "", false
	}

	return pathParts[2], true
}

func HandleTask(w http.ResponseWriter, r *http.Request, service *service.Service) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := validateToken(r, service); err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := service.GetTaskService().CreateTask()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TaskResponse{TaskID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func mapStatusToString(status repo.TaskStatus) string {
	switch status {
	case repo.InProgress:
		return "in_progress"
	case repo.Ready:
		return "ready"
	default:
		return "unknown"
	}
}

func validateToken(r *http.Request, service *service.Service) error {
	token := r.Header.Get("Authorization")

	if _, ok := service.GetUserService().ValidateToken(token); !ok {
		return errors.New("Unauthorized")
	}
	return nil
}
