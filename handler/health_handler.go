package handler

import (
	"encoding/json"
	"net/http"
	"task-manager-api/repository"
	"task-manager-api/usecase"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Cache    string `json:"cache"`
}

func RegisterHealthRoutes(mux *http.ServeMux, repo repository.TaskRepository, cache *usecase.CacheService) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		getHealth(w, r, repo, cache)
	})
}

func getHealth(w http.ResponseWriter, r *http.Request, repo repository.TaskRepository, cache *usecase.CacheService) {
	_, dbErr := repo.GetAll()

	var dbStatus string
	if dbErr != nil {
		dbStatus = "error"
	} else {
		dbStatus = "connected"
	}

	cache.Set("health_check", "ok")

	_, found := cache.Get("health_check")

	var cacheStatus string
	if !found {
		cacheStatus = "error"
	} else {
		cacheStatus = "connected"
	}

	overallStatus := "healthy"

	if dbStatus == "error" || cacheStatus == "error" {
		overallStatus = "unhealthy"
	}

	w.Header().Set("Content-Type", "application/json")

	httpStatus := http.StatusOK
	if overallStatus == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}
	w.WriteHeader(httpStatus)

	response := HealthResponse{
		Status:   overallStatus,
		Database: dbStatus,
		Cache:    cacheStatus,
	}

	json.NewEncoder(w).Encode(response)
}
