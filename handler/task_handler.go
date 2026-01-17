package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-manager-api/domain"
	"task-manager-api/repository"
)

func RegisterTaskRoutes(repo *repository.InMemoryTaskRepository) {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllTasks(w, r, repo)

		case http.MethodPost:
			createTask(w, r, repo)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func getAllTasks(w http.ResponseWriter, r *http.Request, repo *repository.InMemoryTaskRepository) {
	w.Header().Set("Content-Type", "application/json")

	tasks := repo.GetAll()
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "Error marshaling tasks", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(jsonData))
}

func createTask(w http.ResponseWriter, r *http.Request, repo *repository.InMemoryTaskRepository) {
	var newTask domain.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	createdTask, err := json.Marshal(repo.Create(newTask))
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(createdTask))
}
