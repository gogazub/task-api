package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gogazub/app/internal/repo"
	"github.com/gogazub/app/internal/service"
	"github.com/gogazub/app/internal/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Err   string `json:"error"`
}

type TaskResponse struct {
	TaskID string `json:"task_id"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ResultResponse struct {
	Result string `json:"result"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request, userService *service.UserService) {

	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	err := userService.Register(req.Username, req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "User %s created!", req.Username)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, userService *service.UserService) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	token, err := userService.Login(req.Username, req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleTask(w http.ResponseWriter, r *http.Request, taskService *service.TaskService) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := taskService.CreateTask()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TaskResponse{TaskID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func HandleTaskStatus(w http.ResponseWriter, r *http.Request, taskService *service.TaskService) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Task ID is required")
		return
	}
	taskID := pathParts[2]

	status, err := taskService.GetTaskStatus(taskID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	var statusStr string
	switch status {
	case repo.InProgress:
		statusStr = "in_progress"
	case repo.Ready:
		statusStr = "ready"
	default:
		statusStr = "unknown"
	}

	response := StatusResponse{Status: statusStr}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleTaskResult(w http.ResponseWriter, r *http.Request, taskService *service.TaskService) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Task ID is required")
		return
	}
	taskID := pathParts[2]

	result, err := taskService.GetTaskResult(taskID)
	if err != nil {
		if strings.Contains(err.Error(), "not ready") {
			utils.WriteJSONError(w, http.StatusAccepted, err.Error())
		} else {
			utils.WriteJSONError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	response := ResultResponse{Result: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
