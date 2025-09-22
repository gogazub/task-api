package main

import (
	"fmt"
	"net/http"

	"github.com/gogazub/app/internal/api"
	"github.com/gogazub/app/internal/repo"
	"github.com/gogazub/app/internal/service"
)

func main() {
	statusRepo := repo.NewStatusRepo()
	taskService := service.NewTaskService(statusRepo)

	userRepo := repo.NewUserRepo()
	userService := service.NewUserService(userRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		api.HandleRegister(w, r, userService)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		api.HandleLogin(w, r, userService)
	})

	mux.HandleFunc("POST /task", func(w http.ResponseWriter, r *http.Request) {
		api.HandleTask(w, r, taskService)
	})

	mux.HandleFunc("GET /status/", func(w http.ResponseWriter, r *http.Request) {
		api.HandleTaskStatus(w, r, taskService)
	})

	mux.HandleFunc("GET /result/", func(w http.ResponseWriter, r *http.Request) {
		api.HandleTaskResult(w, r, taskService)
	})

	// Root endpoint for health check
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task Service API is running"))
	})

	fmt.Println("Server starting on :8000")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
