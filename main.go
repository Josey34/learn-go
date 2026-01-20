package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"task-manager-api/handler"
	"task-manager-api/middleware"
	"task-manager-api/repository"
	"task-manager-api/usecase"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	repo, err := repository.NewSQLiteTaskRepository("tasks.db")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	userRepo, err := repository.NewSQLiteUserRepository("tasks.db")
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}
	defer userRepo.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello from Task Manager API"}`)
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET environment variable not set")
	}

	uc := usecase.NewTaskUsecase(repo)
	authUc := usecase.NewAuthUsecase(userRepo, jwtSecret)
	handler.SetupRoutes(mux, uc, authUc)

	handler := middleware.Chain(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
	)(mux)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
