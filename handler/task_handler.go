package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskByID(w, r, uc)

		case http.MethodPut:
			updateTask(w, r, uc)

		case http.MethodDelete:
			deleteTask(w, r, uc)

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

func getTaskByID(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := uc.GetByID(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Error marshaling task", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(jsonData))
}

func updateTask(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updateReq dto.UpdateTaskDTO
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	updatedTask, err := uc.UpdateTask(id, updateReq)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(updatedTask)
	if err != nil {
		http.Error(w, "Error marshaling task", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(jsonData))
}

func deleteTask(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = uc.DeleteTask(id)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
