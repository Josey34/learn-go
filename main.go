package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager-api/handler"
	"task-manager-api/repository"
	"task-manager-api/usecase"
)

func main() {
	repo, err := repository.NewSQLiteTaskRepository("tasks.db")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello from Task Manager API"}`)
	})

	uc := usecase.NewTaskUsecase(repo)
	handler.SetupRoutes(uc)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
