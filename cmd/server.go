// @title Orders API
package main

import (
	"fmt"
	"net/http"

	"github.com/gogazub/app/internal/api"
	"github.com/gogazub/app/internal/repo"
	"github.com/gogazub/app/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/gogazub/app/docs"
)

func main() {
	statusRepo := repo.NewStatusRepo()

	userRepo := repo.NewUserRepo()

	service := service.NewService(statusRepo, userRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	authHandler := &api.AuthHandler{service}
	mux.HandleFunc("POST /register", authHandler.HandleRegister)
	mux.HandleFunc("POST /login", authHandler.HandleLogin)

	taskHanler := &api.TaskHandler{Svc: service}
	mux.HandleFunc("POST /task", taskHanler.HandleTask)
	mux.HandleFunc("GET /status/", taskHanler.HandleTaskStatus)
	mux.HandleFunc("GET /result/", taskHanler.HandleTaskResult)

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
