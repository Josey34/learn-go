package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-manager-api/handler"
	"task-manager-api/middleware"
	"task-manager-api/repository"
	"task-manager-api/usecase"
	"time"

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

	cache := usecase.NewCacheService(5 * time.Minute)
	uc := usecase.NewTaskUsecase(repo, cache)
	authUc := usecase.NewAuthUsecase(userRepo, jwtSecret)
	processor := usecase.NewTaskProcessor(repo)

	handler.SetupRoutes(mux, uc, authUc, processor, cache, repo)

	rateLimiter := middleware.NewRateLimiter(20)

	handler := middleware.Chain(
		rateLimiter.RateLimiterMiddleware,
		middleware.CorrelationIDMiddleware,
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
	)(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")

}
