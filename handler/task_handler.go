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

func RegisterTaskRoutes(mux *http.ServeMux, uc *usecase.TaskUsecase) {
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllTasks(w, r, uc)

		case http.MethodPost:
			createTask(w, r, uc)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
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

	tasks, err := uc.GetAllTasks(r.Context())
	if err != nil {
		HandleError(w, err)
		return
	}
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		HandleError(w, err)
		return
	}
	fmt.Fprintf(w, string(jsonData))
}

func createTask(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	var newTask dto.CreateTaskDTO

	err := json.NewDecoder(r.Body).Decode(&newTask)

	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	createdTask, err := uc.CreateTask(r.Context(), newTask)

	if err != nil {
		HandleError(w, err)
		return
	}

	taskResponse, err := json.Marshal(createdTask)
	if err != nil {
		HandleError(w, err)
		return
	}

	fmt.Fprintf(w, string(taskResponse))
}

func getTaskByID(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(w, err)
		return
	}

	task, err := uc.GetByID(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		HandleError(w, err)
		return
	}

	fmt.Fprintf(w, string(jsonData))
}

func updateTask(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(w, err)
		return
	}

	var updateReq dto.UpdateTaskDTO
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	updatedTask, err := uc.UpdateTask(r.Context(), id, updateReq)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonData, err := json.Marshal(updatedTask)
	if err != nil {
		HandleError(w, err)
		return
	}

	fmt.Fprintf(w, string(jsonData))
}

func deleteTask(w http.ResponseWriter, r *http.Request, uc *usecase.TaskUsecase) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/tasks/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(w, err)
		return
	}

	err = uc.DeleteTask(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
