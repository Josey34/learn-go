package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-manager-api/dto"
	"task-manager-api/usecase"
)

func RegisterTaskRoutes(uc *usecase.TaskUsecase) {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllTasks(w, r, uc)

		case http.MethodPost:
			createTask(w, r, uc)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func getAllTasks(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	w.Header().Set("Content-Type", "application/json")

	tasks, err := uc.GetAllTasks()
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "Error marshaling tasks", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(jsonData))
}

func createTask(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	var newTask dto.CreateTaskDTO

	err := json.NewDecoder(r.Body).Decode(&newTask)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	createdTask, err := uc.CreateTask(newTask)

	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	taskResponse, err := json.Marshal(createdTask)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(taskResponse))
}
