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

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
