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
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task Service API is running"))
	})

	fmt.Println("Server starting on :8000")
	go func() {
		var input string
		for {
			fmt.Print("Enter 'help' for API documentation: ")
			fmt.Scan(&input)
			if input == "help" {
				printAPIDocumentationReminder()
			}
		}
	}()

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func printAPIDocumentationReminder() {
	fmt.Print(`
API Endpoints:

1. POST /register
   - Registers a new user.
   - Request body: { "username": "string", "password": "string" }
   - Response: { "token": "string" }

2. POST /login
   - Logs in a user and returns a token.
   - Request body: { "username": "string", "password": "string" }
   - Response: { "token": "string" }

3. POST /task
   - Creates a new task.
   - Response: { "task_id": "string" }

4. GET /status/{task_id}
   - Retrieves the status of a task.
   - Response: { "status": "in_progress" | "ready" | "unknown" }

5. GET /result/{task_id}
   - Retrieves the result of a completed task.
   - Response: { "result": "string" }

6. GET / 
   - Health check for the server.
   - Response: "Task Service API is running"
`)
}
