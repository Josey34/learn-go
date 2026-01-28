package handler

import (
	"net/http"
	"task-manager-api/repository"
	"task-manager-api/usecase"
)

func SetupRoutes(mux *http.ServeMux, uc *usecase.TaskUsecase, authUc *usecase.AuthUsecase, processor *usecase.TaskProcessor, cache *usecase.CacheService, repo repository.TaskRepository) {
	RegisterTaskRoutes(mux, uc)
	RegisterAuthRoutes(mux, authUc)
	RegisterBackgroundRoutes(mux, processor)
	RegisterCacheRoutes(mux, cache)
	RegisterHealthRoutes(mux, repo, cache)
}
