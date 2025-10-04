package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gogazub/app/internal/model"
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

type TaskHandler struct {
	Svc *service.Service
}

type TaskRequest struct {
	Text string `json:"text"`
}

// @Summary Create new task
// @Description Create new task and return task`s id
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header string true "session's token"
// @Param body body model.Task true "code to run"
// @Success 200 {object} TaskResponse
// @Failure 401 {object} ErrorResponse "Unathorized"
// @Failure 404 {object} ErrorResponse ""
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Router /task [post]
func (handler TaskHandler) HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := validateToken(r, handler.Svc); err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := handler.Svc.GetTaskService().CreateTask()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var taskPayload model.Task
	json.NewDecoder(r.Body).Decode(&taskPayload)
	err = handler.Svc.GetProducer().SendMessage(taskPayload)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TaskResponse{TaskID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

// @Summary Return status of task by id
// @Description Return status of task by id
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header string true "session`s token"
// @Param task_id path string true "task_id"
// @Success 200 {object} StatusResponse
// @Failure 401 {object} ErrorResponse "Unathorized"
// @Failure 404 {object} ErrorResponse "Status not found"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Router /status/{id} [get]
func (handler TaskHandler) HandleTaskStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := validateToken(r, handler.Svc); err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	taskID, valid := handleTaskID(w, r)
	if !valid {
		return
	}

	status, err := handler.Svc.GetTaskService().GetTaskStatus(taskID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	statusStr := mapStatusToString(status)

	response := StatusResponse{Status: statusStr}
	writeJSON(w, http.StatusOK, response)
}

// @Summary Return result of task by id
// @Description Return result of task by id
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header string true "session`s token"
// @Param task_id path string true "Task id"
// @Success 200 {object} ResultResponse "return result"
// @Success 202 {object} ResultResponse "task is running"
// @Failure 401 {object} ErrorResponse "Unathorized"
// @Failure 404 {object} ErrorResponse "Result not found"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Router /result/{task_id} [get]
func (handler TaskHandler) HandleTaskResult(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := validateToken(r, handler.Svc); err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	taskID, valid := handleTaskID(w, r)
	if !valid {
		return
	}

	result, err := handler.Svc.GetTaskService().GetTaskResult(taskID)
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
	tokenString := r.Header.Get("Authorization")
	username, err := service.GetUsernameFromToken(tokenString)
	if err != nil {
		return err
	}
	exists := service.GetUserService().FindUser(username)
	if !exists {
		return fmt.Errorf("no such user: %s", username)
	}
	return nil
}
