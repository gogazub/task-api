package api

import (
	"net/http"

	"github.com/gogazub/app/internal/service"
)

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request, service *service.Service) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Username == "" || req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	if err := service.GetUserService().Register(req.Username, req.Password); err != nil {
		writeJSONError(w, http.StatusConflict, "user already exists")
		return
	}

	token, err := service.GetUserService().Login(req.Username, req.Password)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusCreated, LoginResponse{Token: token})
}

func HandleLogin(w http.ResponseWriter, r *http.Request, service *service.Service) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Username == "" || req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	token, err := service.GetUserService().Login(req.Username, req.Password)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{Token: token})
}
