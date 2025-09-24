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

type AuthHandler struct {
	*service.Service
}

// @Summary Register new user and return token
// @Description Register new user and return token
// @Tags auth,user
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "username + password"
// @Success 201 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid JSON"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "failed to create session"
// @Router	/register [post]
func (handler *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
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

	if err := handler.GetUserService().Register(req.Username, req.Password); err != nil {
		writeJSONError(w, http.StatusConflict, "user already exists")
		return
	}

	token, err := handler.GetUserService().Login(req.Username, req.Password)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusCreated, LoginResponse{Token: token})
}

// @Summary login user
// @Description perform authorization and return session`s token if successful
// @Tags auth,user
// @Accept json
// @Produce json
// @Param body body LoginRequest true "username + password"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid JSON"
// @Failure 401 {object} ErrorResponse "Unathorized"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Router /login [post]
func (handler *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
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

	token, err := handler.GetUserService().Login(req.Username, req.Password)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{Token: token})
}
