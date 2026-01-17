package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager-api/handler"
	"task-manager-api/repository"
)

func main() {
	repo := repository.NewInMemoryTaskRepository()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello from Task Manager API"}`)
	})

	handler.SetupRoutes(repo)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
