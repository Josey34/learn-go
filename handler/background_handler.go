package handler

import (
	"encoding/json"
	"net/http"
	"task-manager-api/domain"
	"task-manager-api/dto"
	"task-manager-api/usecase"
)

// RegisterBackgroundRoutes registers background processing routes
func RegisterBackgroundRoutes(mux *http.ServeMux, processor *usecase.TaskProcessor) {
	mux.HandleFunc("POST /tasks/process", func(w http.ResponseWriter, r *http.Request) {
		processTasksInBackground(w, r, processor)
	})
}

// processTasksInBackground processes multiple tasks in the background (fire and forget)
func processTasksInBackground(w http.ResponseWriter, r *http.Request, processor *usecase.TaskProcessor) {
	var req struct {
		TaskIDs []int `json:"task_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, &domain.ValidationError{
			Field:   "body",
			Message: "invalid request body",
		})
		return
	}

	if len(req.TaskIDs) == 0 {
		HandleError(w, &domain.ValidationError{
			Field:   "task_ids",
			Message: "task_ids cannot be empty",
		})
		return
	}

	// Process tasks in background (fire and forget)
	processor.ProcessTasksInBackground(req.TaskIDs)

	// Respond immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(dto.SuccessResponse{
		Message: "Tasks queued for processing",
		Data: map[string]interface{}{
			"count": len(req.TaskIDs),
		},
	})
}
